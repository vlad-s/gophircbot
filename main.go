package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/vlad-s/gophirc"
	"github.com/vlad-s/gophirc/config"
	"github.com/vlad-s/gophircbot/bot"
)

var (
	irc *gophirc.IRC
	log = logrus.New()
)

func init() {
	var count uint8
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			count++
			if count > 1 {
				fmt.Println()
				log.Warnln("Forcefully exiting")
				os.Exit(1)
			}
			fmt.Println()
			log.WithField("reason", "SIGINT").Infoln("Quitting")
			irc.Quit()
		}
	}()
}

func main() {
	conf, err := config.Parse("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = conf.Check()
	if err != nil {
		log.Fatalln(err)
	}

	irc = gophirc.New()
	err = irc.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	bot.AddBasicCallbacks(irc)
	bot.AddCTCPCallbacks(irc)

	if err = irc.Loop(); err != nil {
		log.Fatalln(err)
	}
}
