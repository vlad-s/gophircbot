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
		if e.User.IsAdmin() {
			irc.CTCP(e.User.Nick, e.Code, "OK")
			irc.SendRaw(strings.Join(e.Arguments, " "))
		} else {
			irc.CTCP(e.User.Nick, e.Code, "NOTOK NOT_AN_ADMIN")
		}
	}).AddEventCallback("QUIT", func(e *gophirc.Event) {
		if e.User.IsAdmin() {
			irc.CTCP(e.User.Nick, e.Code, "OK")
			irc.Quit()
		} else {
			irc.CTCP(e.User.Nick, e.Code, "NOTOK NOT_AN_ADMIN")
		}
	})
}

func AddBasicCallbacks(irc *gophirc.IRC) {
	irc.AddEventCallback("PRIVMSG", func(e *gophirc.Event) {
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
					return // todo: do something with the error?
				}
				irc.PrivMsg(replyTo, fmt.Sprintf("[URL] %s", title))
			case v == "shrug":
				irc.PrivMsg(replyTo, `¯\_(ツ)_/¯`)
			}
		}
	})
}
