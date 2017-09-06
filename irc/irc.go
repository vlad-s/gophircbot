package irc

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
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
	Connected chan struct{}

	registered bool
	Registered chan struct{}

	Identified chan struct{}
}

type Event struct {
	Raw       string
	Code      string
	Source    string
	User      *User
	Arguments []string
	Message   string
}

type IRC struct {
	conn   net.Conn
	Server Server

	State  State
	Events map[string][]func(*Event)

	raw  chan string
	quit chan struct{}
}

func (irc *IRC) Connect() error {
	dest := fmt.Sprintf("%s:%d", irc.Server.Address, irc.Server.Port)
	c, err := net.DialTimeout("tcp", dest, 5*time.Second)
	if err != nil {
		return errors.Wrap(err, "Error dialing the host")
	}
	irc.conn = c
	return nil
}

func (irc *IRC) Disconnect() {
	fmt.Fprint(irc.conn, "QUIT :Quitting\r\n")
	irc.conn.Close()
}

func (irc *IRC) Loop() error {
	var gracefulExit bool

	go func() {
		for {
			select {
			case <-irc.quit:
				gracefulExit = true
				irc.Disconnect()
				return
			case s := <-irc.raw:
				logger.Log.Debugf("Raw:\t%q", s)
			}
		}
	}()

	s := bufio.NewScanner(irc.conn)
	for s.Scan() {
		go irc.getRaw(s.Text())
	}

	if gracefulExit {
		return nil
	}
	return errors.Wrap(s.Err(), "Error while looping")
}

func (irc *IRC) AddEventCallback(code string, cb func(*Event)) *IRC {
	irc.Events[code] = append(irc.Events[code], cb)
	return irc
}

func (irc *IRC) parseToEvent(raw string) (event *Event, ok bool) {
	irc.raw <- raw
	event = &Event{Raw: raw}
	if raw[0] != ':' {
		return
	}

	raw = raw[1:]
	split := strings.Split(raw, " ")

	event.Source = split[0]
	event.Code = split[1]
	event.Arguments = split[2:]

	if u, ok := ParseUser(event.Source); ok {
		event.User = u
	}

	if event.Code == "PRIVMSG" {
		message := strings.Join(event.Arguments[1:], " ")[1:]
		if IsCTCP(message) {
			message = strings.Trim(message, "\001")
			message_args := strings.Split(message, " ")
			event.Code = message_args[0]
			event.Arguments = message_args[1:]
		}
	}

	return event, true
}

func (irc *IRC) getRaw(raw string) {
	e, ok := irc.parseToEvent(raw)
	if !ok {
		split := strings.Split(e.Raw, " ")
		if split[0] == "PING" {
			irc.pong(split[1])
		}
	}

	for k, v := range irc.Events {
		if k != e.Code {
			continue
		}
		fmt.Println("calling funcs for k", k)
		for _, f := range v {
			f(e)
		}
	}

	switch e.Code {
	case "404":
		logger.Log.WithField("channel", e.Arguments[0]).Warnln("Can't send to channel")
	case "474":
		logger.Log.WithField("channel", e.Arguments[0]).Warnln("Can't join channel")
	case "KICK":
		if e.Arguments[1] == config.Get().Nickname {
			logger.Log.WithFields(logger.Fields(map[string]interface{}{
				"user": e.User.Nick, "channel": e.Arguments[0],
			})).Warnln("We got kicked from a channel")
		}
	}
}

func (irc *IRC) parseCommand(replyTo, message string, user *User) {
	message = strings.TrimPrefix(message, ",")
	s := strings.Split(message, " ")

	switch s[0] {
	case "help":
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, %s", user.Nick, HELP))
	case "say":
		irc.PrivMsg(replyTo, strings.Join(s[1:], " "))
	case "yell":
		irc.PrivMsg(replyTo, strings.ToUpper(strings.Join(s[1:], " ")))
	case "slap":
		if len(s) < 2 {
			irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,slap user", user.Nick))
			return
		}
		irc.Action(replyTo, fmt.Sprintf("slaps %s around a bit with a large trout", s[1]))
	case "join":
		if !user.IsAdmin() {
			irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, can't let you do that.", user.Nick))
			return
		}
		if len(s) < 2 || s[1][0] != '#' {
			irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,join #channel", user.Nick))
			return
		}
		irc.Join(s[1])
	case "part":
		if !user.IsAdmin() {
			irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, can't let you do that.", user.Nick))
			return
		}
		if len(s) < 2 || s[1][0] != '#' {
			irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,part #channel", user.Nick))
			return
		}
		irc.Part(s[1])
	}
}

func (irc *IRC) AddCTCPCallbacks() {
	irc.AddEventCallback("VERSION", func(e *Event) {
		irc.Notice(e.User.Nick, fmt.Sprintf("\001VERSION %s\001", VERSION))
	}).AddEventCallback("TIME", func(e *Event) {
		irc.Notice(e.User.Nick, fmt.Sprintf("\001TIME %s\001", time.Now().Format(time.RFC850)))
	}).AddEventCallback("PING", func(e *Event) {
		irc.Notice(e.User.Nick, fmt.Sprintf("\001PING %s\001", e.Arguments[0]))
	}).AddEventCallback("RAW", func(e *Event) {
		if e.User.IsAdmin() {
			irc.CTCP(e.User.Nick, e.Code, "OK")
			irc.sendRaw(strings.Join(e.Arguments, " "))
		} else {
			irc.CTCP(e.User.Nick, e.Code, "NOTOK NOT_AN_ADMIN")
		}
	}).AddEventCallback("QUIT", func(e *Event) {
		if e.User.IsAdmin() {
			irc.CTCP(e.User.Nick, e.Code, "OK")
			irc.Quit()
		} else {
			irc.CTCP(e.User.Nick, e.Code, "NOTOK NOT_AN_ADMIN")
		}
	})
}

func (irc *IRC) Quit() {
	irc.quit <- struct{}{}
}

func New(s Server) *IRC {
	i := &IRC{
		Server: s,

		State: State{
			Connected:  make(chan struct{}),
			Registered: make(chan struct{}),
			Identified: make(chan struct{}),
		},
		Events: make(map[string][]func(*Event)),

		raw:  make(chan string),
		quit: make(chan struct{}, 1),
	}

	i.AddEventCallback("NOTICE", func(e *Event) {
		if e.Arguments[0] == "*" && !i.State.registered {
			i.State.Connected <- struct{}{}
			i.Register()
		}
	}).AddEventCallback("001", func(e *Event) {
		i.Identify()
	}).AddEventCallback("900", func(e *Event) {
		i.State.Identified <- struct{}{}
		for _, v := range config.Get().Server.Channels {
			i.Join(v)
		}
	}).AddEventCallback("INVITE", func(e *Event) {
		channel := e.Arguments[1][1:]
		i.Join(channel)
		i.PrivMsg(channel, fmt.Sprintf("Hi %s, %s invited me here.", channel, e.User.Nick))
	})

	i.AddEventCallback("PRIVMSG", func(e *Event) {
		replyTo := e.Arguments[0]
		message := strings.Join(e.Arguments[1:], " ")[1:]

		for _, v := range strings.Split(message, " ") {
			switch true {
			case IsValidURL(v):
				if ok, _ := regexp.MatchString(`https?://(www\.)?(filelist\.ro|flro\.org)`, v); ok {
					return
				}

				title, err := GetTitle(v)
				if err != nil {
					if config.Get().Debug {
						logger.Log.Errorln(v, err)
					}
					return // todo: do something with the error?
				}
				i.PrivMsg(replyTo, fmt.Sprintf("[URL] %s", title))
			case v == "shrug":
				i.PrivMsg(replyTo, `¯\_(ツ)_/¯`)
			}
		}
	})

	i.AddCTCPCallbacks()

	return i
}
