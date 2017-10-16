package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/vlad-s/gophirc"
	"github.com/vlad-s/gophirc/config"
	"github.com/vlad-s/gophirc/logger"
	"github.com/vlad-s/gophircbot/apiconfig"
	"github.com/vlad-s/gophircbot/bot"
	"github.com/vlad-s/gophircbot/db"
)

var (
	servers []*gophirc.IRC

	configFlag    = flag.String("config", "config.json", "Path to the config `file`")
	apiConfigFlag = flag.String("api_config", "config_api.json", "Path to the API config `file`")
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
				logger.Log.Warnln("Forcefully exiting")
				os.Exit(1)
			}
			fmt.Println()
			logger.Log.WithField("reason", "SIGINT").Infoln("Quitting")

			for _, v := range servers {
				v.Quit()
			}
		}
	}()
}

func main() {
	flag.Parse()

	conf, err := config.Parse(*configFlag)
	if err != nil {
		logger.Log.Fatalln(err)
	}

	err = conf.Check()
	if err != nil {
		logger.Log.Fatalln(err)
	}

	apiConf, err := apiconfig.Parse(*apiConfigFlag)
	if err != nil {
		logger.Log.Fatalln(err)
	}

	botdb, err := db.Connect(apiConf.Database)
	if err != nil {
		logger.Log.Fatalln(err)
	}
	defer botdb.Close()

	logger.Log.Infoln("Connected to database")
	db.AutoMigrate(botdb)

	var wg sync.WaitGroup
	for _, server := range conf.Servers {
		irc := gophirc.New(server, &wg)
		err = irc.Connect()
		if err != nil {
			logger.Log.Fatalln(err)
		}

		servers = append(servers, irc)

		bot.AddBasicCallbacks(irc)
		bot.AddCTCPCallbacks(irc)

		go irc.Loop()
	}
	wg.Wait()
}
