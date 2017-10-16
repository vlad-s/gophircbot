package bot

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/vlad-s/gophirc"
)

// AddCTCPCallbacks adds the CTCP set of callbacks
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
	}).AddEventCallback("QUITBOT", func(e *gophirc.Event) {
		if irc.IsAdmin(e.User) {
			irc.CTCP(e.User.Nick, e.Code, "OK")
			irc.Quit()
		} else {
			irc.CTCP(e.User.Nick, e.Code, "NOTOK NOT_AN_ADMIN")
		}
	})
}

// AddBasicCallbacks adds the basic set of callbacks.
func AddBasicCallbacks(irc *gophirc.IRC) {
	irc.AddEventCallback("PRIVMSG", func(e *gophirc.Event) {
		if IsIgnored(e.User.Nick) {
			return
		}

		message := strings.Join(e.Arguments[1:], " ")[1:]
		for _, v := range strings.Split(message, " ") {
			if !IsValidURL(v) {
				continue
			}
			if ok, _ := regexp.MatchString(`https?://(www\.)?(filelist\.ro|flro\.org)`, v); ok {
				return
			}
			title, err := GetTitle(v)
			if err != nil {
				return
			}
			irc.PrivMsgf(e.ReplyTo, "[URL] %s", title)
		}
	}).AddEventCallback("PRIVMSG", func(e *gophirc.Event) {
		if IsIgnored(e.User.Nick) {
			return
		}

		message := strings.Join(e.Arguments[1:], " ")[1:]

		addStaticReplies(message, irc, e)
		addCommands(message, irc, e)
		addAdminCommands(message, irc, e)
	})
}

// addCommands adds public available commands.
func addCommands(message string, irc *gophirc.IRC, e *gophirc.Event) {
	if message[0] == ',' {
		split := strings.Split(message[1:], " ")

		switch split[0] {
		case "say":
			say(irc, e)
		case "yell":
			yell(irc, e)
		case "gif":
			searchGif(irc, e)
		}
	}
}

// addAdminCommands adds admin-only commands.
func addAdminCommands(message string, irc *gophirc.IRC, e *gophirc.Event) {
	if message[0] == ',' && irc.IsAdmin(e.User) {
		split := strings.Split(message[1:], " ")

		switch split[0] {
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
			kickBan(irc, e)
		case "ignore":
			ignoreUser(irc, e)
		case "recognize":
			recognizeUser(irc, e)
		}
	}
}

// addStaticReplies adds static replies, as the name implies.
func addStaticReplies(message string, irc *gophirc.IRC, e *gophirc.Event) {
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
}
