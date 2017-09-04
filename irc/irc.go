package irc

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vlad-s/gophircbot/config"
	"github.com/vlad-s/gophircbot/logger"
)

const (
	VERSION = "gophircbot v0.1, https://github.com/vlad-s/gophircbot"
	HELP    = "gophircbot v0.1, please check https://github.com/vlad-s/gophircbot for more info"
)

type Server struct {
	Address string
	Port    uint16
}

type State struct {
	Connected  bool
	Registered bool
	Identified bool
}

type IRC struct {
	Conn   net.Conn
	Server Server

	State State

	Receiver chan string
	Sent     chan string

	quit chan struct{}
}

func (irc *IRC) Connect() error {
	dest := fmt.Sprintf("%s:%d", irc.Server.Address, irc.Server.Port)
	c, err := net.DialTimeout("tcp", dest, 5*time.Second)
	if err != nil {
		return errors.Wrap(err, "Error dialing the host")
	}
	irc.Conn = c
	return nil
}

func (irc *IRC) Disconnect() {
	irc.Conn.Close()
	irc.State.Connected = false

	close(irc.Receiver)
	close(irc.Sent)
	close(irc.quit)
}

func (irc *IRC) Loop() error {
	var wg sync.WaitGroup
	var gracefulExit bool

	go irc.parseEvent()

	go func() {
		for {
			select {
			case <-irc.quit:
				gracefulExit = true
				irc.sendRaw("QUIT :Quitting")
				irc.Disconnect()
				return
			case s := <-irc.Sent:
				if config.Get().Debug {
					logger.Log.Debugf("Sent:\t%q", s)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		s := bufio.NewScanner(irc.Conn)
		for s.Scan() {
			irc.Receiver <- s.Text()
		}
		wg.Done()
	}()
	wg.Wait()

	if gracefulExit {
		return nil
	}
	return errors.New("Error while looping")
}

func (irc *IRC) sendRaw(s string) {
	fmt.Fprint(irc.Conn, s+"\r\n")
	irc.Sent <- s
}

func (irc *IRC) parseEvent() {
	for s := range irc.Receiver {
		if config.Get().Debug {
			logger.Log.Debugln("parseEvent():\t", s)
		}

		split := strings.Split(s, " ")

		if split[1] == "NOTICE" && split[2] == "*" && !irc.State.Registered {
			irc.State.Connected = true
			go irc.register()
		}

		if split[0] == "PING" {
			go irc.pong(split[1])
		}

		switch split[1] {
		case "001":
			go irc.identify()
		case "900":
			irc.State.Identified = true
			for _, v := range config.Get().Server.Channels {
				go irc.join(v)
			}
		case "404":
			logger.Log.WithField("channel", split[3]).Warnln("Can't send to channel")
		case "474":
			logger.Log.WithField("channel", split[3]).Warnln("Can't join channel")
		case "INVITE":
			user := ParseUser(split[0])
			channel := split[3][1:]
			irc.join(channel)
			irc.privMsg(channel, fmt.Sprintf("Hi %s, %s invited me here.", channel, user.Nick))
		case "KICK":
			if split[3] == config.Get().Nickname {
				user := ParseUser(split[0])
				channel := split[2]

				logger.Log.WithFields(logger.Fields(map[string]interface{}{
					"user": user.Nick, "channel": channel,
				})).Warnln("We got kicked from a channel")
			}
		case "PRIVMSG":
			go irc.parsePrivMsg(split)
		}
	}
}

func (irc *IRC) parsePrivMsg(s []string) {
	user, replyTo := ParseUser(s[0]), s[2]
	if replyTo == config.Get().Nickname {
		replyTo = user.Nick
	}

	message := strings.Join(s[3:], " ")
	message = strings.TrimPrefix(message, ":")
	message = strings.TrimSpace(message)

	switch message[0] {
	case ',':
		irc.parseCommand(replyTo, message, user)
		return
	case '\001':
		irc.parseCtcp(replyTo, message, user)
		return
	}

	s = strings.Split(message, " ")
	for _, v := range s {
		switch true {
		case IsValidURL(v):
			title, err := GetTitle(v)
			if err != nil {
				if config.Get().Debug {
					logger.Log.Errorln(v, err)
				}
				return // todo: do something with the error?
			}
			irc.privMsg(replyTo, fmt.Sprintf("[URL] %s", title))
		case v == "shrug":
			irc.privMsg(replyTo, `¯\_(ツ)_/¯`)
		}
	}
}

func (irc *IRC) parseCommand(replyTo, message string, user User) {
	message = strings.TrimPrefix(message, ",")
	s := strings.Split(message, " ")

	switch s[0] {
	case "help":
		irc.privMsg(replyTo, fmt.Sprintf("%s, %s", user.Nick, HELP))
	case "say":
		irc.privMsg(replyTo, strings.Join(s[1:], " "))
	case "yell":
		irc.privMsg(replyTo, strings.ToUpper(strings.Join(s[1:], " ")))
	case "join":
		if !user.IsAdmin() {
			irc.privMsg(replyTo, fmt.Sprintf("Sorry %s, can't let you do that.", user.Nick))
			return
		}
		if len(s) < 2 || s[1][0] != '#' {
			irc.privMsg(replyTo, fmt.Sprintf("%s, usage: ,join #channel", user.Nick))
			return
		}
		irc.join(s[1])
	case "part":
		if !user.IsAdmin() {
			irc.privMsg(replyTo, fmt.Sprintf("Sorry %s, can't let you do that.", user.Nick))
			return
		}
		if len(s) < 2 || s[1][0] != '#' {
			irc.privMsg(replyTo, fmt.Sprintf("%s, usage: ,part #channel", user.Nick))
			return
		}
		irc.part(s[1])
	}
}

func (irc *IRC) parseCtcp(replyTo, message string, user User) {
	if message[len(message)-1:][0] != '\001' {
		return
	}

	message = strings.Trim(message, "\001")
	s := strings.Split(message, " ")

	switch s[0] {
	case "VERSION":
		irc.notice(replyTo, fmt.Sprintf("\001VERSION %s\001", VERSION))
	case "TIME":
		irc.notice(replyTo, fmt.Sprintf("\001TIME %s\001", time.Now().Format(time.RFC850)))
	case "PING":
		irc.notice(replyTo, fmt.Sprintf("\001PING %s\001", s[1]))
	case "QUIT":
		if user.IsAdmin() {
			irc.ctcp(replyTo, s[0], "OK")
			irc.quit <- struct{}{}
		} else {
			irc.ctcp(replyTo, s[0], "NOTOK NOT_AN_ADMIN")
		}
	case "RAW":
		if user.IsAdmin() {
			irc.ctcp(replyTo, s[0], "OK")
			irc.sendRaw(strings.Join(s[1:], " "))
		} else {
			irc.ctcp(replyTo, s[0], "NOTOK NOT_AN_ADMIN")
		}
	}

}

func (irc *IRC) pong(s string) {
	irc.sendRaw(fmt.Sprintf("PONG %s", s))
}

func (irc *IRC) register() {
	irc.sendRaw(fmt.Sprintf("USER %s 8 * %s", config.Get().Username, config.Get().Realname))
	irc.sendRaw(fmt.Sprintf("NICK %s", config.Get().Nickname))

	irc.State.Registered = true
}

func (irc *IRC) identify() {
	if config.Get().Server.NickservPassword == "" {
		return
	}
	irc.sendRaw(fmt.Sprintf("NS identify %s", config.Get().Server.NickservPassword))
}

func (irc *IRC) join(channel string) {
	irc.sendRaw(fmt.Sprintf("JOIN %s", channel))
}

func (irc *IRC) part(channel string) {
	irc.sendRaw(fmt.Sprintf("PART %s", channel))
}

func (irc *IRC) privMsg(replyTo, message string) {
	irc.sendRaw(fmt.Sprintf("PRIVMSG %s :%s", replyTo, message))
}

func (irc *IRC) notice(replyTo, message string) {
	irc.sendRaw(fmt.Sprintf("NOTICE %s :%s", replyTo, message))
}

func (irc *IRC) action(replyTo, message string) {
	irc.sendRaw(fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001", replyTo, message))
}

func (irc *IRC) ctcp(replyTo, ctcp, message string) {
	irc.notice(replyTo, fmt.Sprintf("\001%s %s\001", ctcp, message))
}

func (irc *IRC) Quit(s string) {
	irc.quit <- struct{}{}
}

func New(s Server) *IRC {
	return &IRC{
		Server: s,

		Receiver: make(chan string),
		Sent:     make(chan string),

		quit: make(chan struct{}, 1),
	}
}
