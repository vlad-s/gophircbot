package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/vlad-s/gophircbot/config"
	"github.com/vlad-s/gophircbot/irc"
	"github.com/vlad-s/gophircbot/logger"
)

var bot *irc.IRC

func init() {
	var count uint8
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			count++
			if count > 1 {
				fmt.Println()
				logger.Log.Warnln("Forcefully exiting")
				os.Exit(1)
			}
			fmt.Println()
			logger.Log.WithField("reason", "SIGINT").Infoln("Quitting")
			bot.Quit("SIGINT")
		}
	}()
}

func main() {
	logger.Log.Infoln("Starting up")

	logger.Log.Infoln("Reading config.json")
	conf, err := config.Parse("config.json")
	if err != nil {
		logger.Log.Fatalln(err)
	}

	logger.Log.Infoln("Checking the config for errors")
	err = conf.Check()
	if err != nil {
		logger.Log.Fatalln(err)
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"server": conf.Server.Address, "port": conf.Server.Port,
	})).Infoln("Connecting to server")

	bot = irc.New(irc.Server{
		Address: conf.Server.Address,
		Port:    conf.Server.Port,
	})

	err = bot.Connect()
	if err != nil {
		logger.Log.Fatalln(err)
	}
	defer bot.Disconnect()

	go logStates()

	if err = bot.Loop(); err != nil {
		logger.Log.Fatalln(err)
	}

	logger.Log.Infoln("Exiting")
}

func logStates() {
	var c, r, i bool
	for {
		if bot.State.Connected && !c {
			logger.Log.Infoln("Successfully connected to server")
			c = true
		}
		if bot.State.Registered && !r {
			logger.Log.Infoln("Successfully registered on network")
			r = true
		}
		if bot.State.Identified && !i {
			logger.Log.Infoln("Successfully identified to Nickserv")
			i = true
		}

		if c && r && i {
			break
		}
	}
}
