package irc

import (
	"fmt"

	"github.com/vlad-s/gophircbot/config"
)

func (irc *IRC) sendRaw(s string) {
	irc.raw <- s
	fmt.Fprint(irc.conn, s+"\r\n")
}

func (irc *IRC) pong(s string) {
	irc.sendRaw(fmt.Sprintf("PONG %s", s))
}

func (irc *IRC) Register() {
	irc.sendRaw(fmt.Sprintf("USER %s 8 * %s", config.Get().Username, config.Get().Realname))
	irc.sendRaw(fmt.Sprintf("NICK %s", config.Get().Nickname))

	irc.State.registered = true
	irc.State.Registered <- struct{}{}
}

func (irc *IRC) Identify() {
	if config.Get().Server.NickservPassword == "" {
		return
	}
	irc.sendRaw(fmt.Sprintf("NS identify %s", config.Get().Server.NickservPassword))
}

func (irc *IRC) Join(channel string) {
	irc.sendRaw(fmt.Sprintf("JOIN %s", channel))
}

func (irc *IRC) Part(channel string) {
	irc.sendRaw(fmt.Sprintf("PART %s", channel))
}

func (irc *IRC) PrivMsg(replyTo, message string) {
	irc.sendRaw(fmt.Sprintf("PRIVMSG %s :%s", replyTo, message))
}

func (irc *IRC) Notice(replyTo, message string) {
	irc.sendRaw(fmt.Sprintf("NOTICE %s :%s", replyTo, message))
}

func (irc *IRC) Action(replyTo, message string) {
	irc.sendRaw(fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001", replyTo, message))
}

func (irc *IRC) CTCP(replyTo, ctcp, message string) {
	irc.Notice(replyTo, fmt.Sprintf("\001%s %s\001", ctcp, message))
}

func (irc *IRC) Kick(channel, nick, message string) {
	if message != "" {
		message = ":" + message
	}
	irc.sendRaw(fmt.Sprintf("KICK %s %s %s", channel, nick, message))
}

func (irc *IRC) Invite(nick, channel string) {
	irc.sendRaw(fmt.Sprintf("INVITE %s %s", nick, channel))
}

func (irc *IRC) Mode(channel, mode, nick string) {
	irc.sendRaw(fmt.Sprintf("MODE %s %s %s", channel, mode, nick))
}

func (irc *IRC) Ban(channel, nick string) {
	irc.Mode(channel, "+b", nick)
}

func (irc *IRC) Unban(channel, nick string) {
	irc.Mode(channel, "-b", nick)
}

func (irc *IRC) KickBan(channel, nick string) {
	irc.Ban(channel, nick)
	irc.Kick(channel, nick, "beep boop i press buttons")
}
