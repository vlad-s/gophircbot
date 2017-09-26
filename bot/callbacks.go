package bot

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/vlad-s/gophirc"
)

func AddCTCPCallbacks(irc *gophirc.IRC) {
	irc.AddEventCallback("VERSION", func(e *gophirc.Event) {
		irc.Notice(e.User.Nick, fmt.Sprintf("\001VERSION %s\001", VERSION))
	}).AddEventCallback("TIME", func(e *gophirc.Event) {
		irc.Notice(e.User.Nick, fmt.Sprintf("\001TIME %s\001", time.Now().Format(time.RFC850)))
	}).AddEventCallback("PING", func(e *gophirc.Event) {
		irc.Notice(e.User.Nick, fmt.Sprintf("\001PING %s\001", e.Arguments[0]))
	}).AddEventCallback("RAW", func(e *gophirc.Event) {
		if irc.IsAdmin(e.User) {
			irc.CTCP(e.User.Nick, e.Code, "OK")
			irc.SendRaw(strings.Join(e.Arguments, " "))
		} else {
			irc.CTCP(e.User.Nick, e.Code, "NOTOK NOT_AN_ADMIN")
		}
	}).AddEventCallback("QUIT", func(e *gophirc.Event) {
		if irc.IsAdmin(e.User) {
			irc.CTCP(e.User.Nick, e.Code, "OK")
			irc.Quit()
		} else {
			irc.CTCP(e.User.Nick, e.Code, "NOTOK NOT_AN_ADMIN")
		}
	})
}

func AddBasicCallbacks(irc *gophirc.IRC) {
	irc.AddEventCallback("PRIVMSG", func(e *gophirc.Event) {
		message := strings.Join(e.Arguments[1:], " ")[1:]

		for _, v := range strings.Split(message, " ") {
			if IsValidURL(v) {
				if ok, _ := regexp.MatchString(`https?://(www\.)?(filelist\.ro|flro\.org)`, v); ok {
					return
				}
				title, err := GetTitle(v)
				if err != nil {
					return
				}
				irc.PrivMsgf(e.ReplyTo, "[URL] %s", title)
			}
		}
	}).AddEventCallback("PRIVMSG", func(e *gophirc.Event) {
		message := strings.Join(e.Arguments[1:], " ")[1:]

		switch message {
		case "test":
			irc.PrivMsg(e.ReplyTo, "test")
		case "ping":
			irc.PrivMsgf(e.ReplyTo, "pong %s", e.User.Nick)
		case "shrug":
			irc.PrivMsg(e.ReplyTo, `¯\_(ツ)_/¯`)
		case `\o`:
			irc.PrivMsg(e.ReplyTo, `o/`)
		case `o/`:
			irc.PrivMsg(e.ReplyTo, `\o`)
		}

		if ok, _ := regexp.MatchString(`^[Ss]alut\s*[!.]?$`, message); ok {
			irc.PrivMsgf(e.ReplyTo, "Salut, %s!", e.User.Nick)
		}

		if ok, _ := regexp.MatchString(`^[Hh](i|ello)\s*[!.]?$`, message); ok {
			irc.PrivMsgf(e.ReplyTo, "Hello, %s!", e.User.Nick)
		}

		if message[0] == ',' {
			split := strings.Split(message[1:], " ")

			switch split[0] {
			case "say":
				say(irc, e)
			case "yell":
				yell(irc, e)
			case "nick":
				nick(irc, e)
			case "join":
				join(irc, e)
			case "part":
				part(irc, e)
			case "invite":
				invite(irc, e)
			case "k", "kick":
				kick(irc, e)
			case "b", "ban":
				ban(irc, e)
			case "ub", "unban":
				unban(irc, e)
			case "kb", "kickban":
				kickban(irc, e)
			case "gif":
				searchGif(irc, e)
			}
		}
	})
}
