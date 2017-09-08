package bot

import (
	"fmt"
	"strings"

	"github.com/vlad-s/gophirc"
)

func say(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")
	args := strings.Join(split[1:], " ")
	args = strings.TrimSpace(args)

	if len(split) < 2 {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,say text", event.User.Nick))
		return
	}

	irc.PrivMsg(replyTo, args)
}

func yell(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")
	args := strings.Join(split[1:], " ")
	args = strings.TrimSpace(args)

	if len(split) < 2 {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,yell text", event.User.Nick))
		return
	}

	irc.PrivMsg(replyTo, strings.ToUpper(args))
}

func nick(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,nick nickname", event.User.Nick))
		return
	}

	if !irc.IsAdmin(event.User) {
		irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, I can't let you do that.", event.User.Nick))
		return
	}

	irc.Nick(split[1])
}

func join(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")

	if len(split) < 2 || split[1][0] != '#' {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,join #channel", event.User.Nick))
		return
	}

	if !irc.IsAdmin(event.User) {
		irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, I can't let you do that.", event.User.Nick))
		return
	}

	irc.Join(split[1])
}

func part(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")

	if len(split) < 2 || split[1][0] != '#' {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,part #channel", event.User.Nick))
		return
	}

	if !irc.IsAdmin(event.User) {
		irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, I can't let you do that.", event.User.Nick))
		return
	}

	irc.Part(split[1])
}

func invite(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,invite user", event.User.Nick))
		return
	}

	if !irc.IsAdmin(event.User) {
		irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, I can't let you do that.", event.User.Nick))
		return
	}

	irc.Invite(split[1], replyTo)
}

func kick(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,kick user <message>", event.User.Nick))
		return
	}

	if !irc.IsAdmin(event.User) {
		irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, I can't let you do that.", event.User.Nick))
		return
	}

	irc.Kick(replyTo, split[1], strings.Join(split[2:], " "))
}

func ban(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,ban user", event.User.Nick))
		return
	}

	if !irc.IsAdmin(event.User) {
		irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, I can't let you do that.", event.User.Nick))
		return
	}

	irc.Ban(replyTo, split[1])
}

func unban(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,unban user", event.User.Nick))
		return
	}

	if !irc.IsAdmin(event.User) {
		irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, I can't let you do that.", event.User.Nick))
		return
	}

	irc.Unban(replyTo, split[1])
}

func kickban(irc *gophirc.IRC, event *gophirc.Event) {
	replyTo := event.Arguments[0]
	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsg(replyTo, fmt.Sprintf("%s, usage: ,kickban user", event.User.Nick))
		return
	}

	if !irc.IsAdmin(event.User) {
		irc.PrivMsg(replyTo, fmt.Sprintf("Sorry %s, I can't let you do that.", event.User.Nick))
		return
	}

	irc.KickBan(replyTo, split[1])
}
