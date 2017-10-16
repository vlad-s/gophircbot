package bot

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/vlad-s/gophirc"
	"github.com/vlad-s/gophirc/logger"
	"github.com/vlad-s/gophircbot/db"
)

func say(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")
	args := strings.Join(split[1:], " ")
	args = strings.TrimSpace(args)

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,say text", event.User.Nick)
		return
	}

	irc.PrivMsg(event.ReplyTo, args)
}

func yell(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")
	args := strings.Join(split[1:], " ")
	args = strings.TrimSpace(args)

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,yell text", event.User.Nick)
		return
	}

	irc.PrivMsg(event.ReplyTo, strings.ToUpper(args))
}

func nick(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,nick nickname", event.User.Nick)
		return
	}

	irc.Nick(split[1])
}

func join(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")

	if len(split) < 2 || split[1][0] != '#' {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,join #channel", event.User.Nick)
		return
	}

	irc.Join(split[1])
}

func part(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")

	if len(split) < 2 || split[1][0] != '#' {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,part #channel", event.User.Nick)
		return
	}

	irc.Part(split[1])
}

func invite(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,invite user", event.User.Nick)
		return
	}

	irc.Invite(split[1], event.ReplyTo)
}

func kick(irc *gophirc.IRC, event *gophirc.Event) {
	if !gophirc.IsChannel(event.ReplyTo) {
		irc.PrivMsgf(event.ReplyTo, "%s, command must be sent on a channel.", event.User.Nick)
		return
	}

	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,kick user <message>", event.User.Nick)
		return
	}

	irc.Kick(event.ReplyTo, split[1], strings.Join(split[2:], " "))
}

func ban(irc *gophirc.IRC, event *gophirc.Event) {
	if !gophirc.IsChannel(event.ReplyTo) {
		irc.PrivMsgf(event.ReplyTo, "%s, command must be sent on a channel.", event.User.Nick)
		return
	}

	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,ban user", event.User.Nick)
		return
	}

	irc.Ban(event.ReplyTo, split[1])
}

func unban(irc *gophirc.IRC, event *gophirc.Event) {
	if !gophirc.IsChannel(event.ReplyTo) {
		irc.PrivMsgf(event.ReplyTo, "%s, command must be sent on a channel.", event.User.Nick)
		return
	}

	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,unban user", event.User.Nick)
		return
	}

	irc.Unban(event.ReplyTo, split[1])
}

func kickBan(irc *gophirc.IRC, event *gophirc.Event) {
	if !gophirc.IsChannel(event.ReplyTo) {
		irc.PrivMsgf(event.ReplyTo, "%s, command must be sent on a channel.", event.User.Nick)
		return
	}

	split := strings.Split(event.Message, " ")

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,kickban user", event.User.Nick)
		return
	}

	irc.KickBan(event.ReplyTo, split[1])
}

func searchGif(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")
	query := strings.Join(split[1:], " ")

	reply, err := GetGif(query)
	if err != nil {
		return
	}

	irc.PrivMsgf(event.ReplyTo, reply)
}

func ignoreUser(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")
	nick := split[1]

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,ignore user", event.User.Nick)
		return
	}

	if IsIgnored(nick) {
		irc.PrivMsgf(event.ReplyTo, "User %q is already ignored.", nick)
		return
	}

	if err := db.Get().Create(&db.IgnoredUser{Nick: nick}).Error; err != nil {
		logger.Log.WithField("nick", nick).
			Errorln(errors.Wrap(err, "Error inserting ignored nick in database"))
	} else {
		irc.PrivMsgf(event.ReplyTo, "User %q successfully ignored.", nick)
	}
}

func recognizeUser(irc *gophirc.IRC, event *gophirc.Event) {
	split := strings.Split(event.Message, " ")
	nick := split[1]

	if len(split) < 2 {
		irc.PrivMsgf(event.ReplyTo, "%s, usage: ,recognize user", event.User.Nick)
		return
	}

	if !IsIgnored(nick) {
		irc.PrivMsgf(event.ReplyTo, "User %q is not ignored.", nick)
		return
	}

	database := db.Get()
	var ignoredUser db.IgnoredUser
	database.Where("nick = ?", nick).First(&ignoredUser)
	if err := database.Delete(&ignoredUser).Error; err != nil {
		logger.Log.WithField("nick", nick).
			Errorln(errors.Wrap(err, "Error deleting ignored user from the database"))
	} else {
		irc.PrivMsgf(event.ReplyTo, "User %q will be recognized from now on")
	}
}
