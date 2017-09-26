package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/vlad-s/gophirc"
	"github.com/vlad-s/gophirc/config"
	"github.com/vlad-s/gophircbot/api_config"
	"github.com/vlad-s/gophircbot/bot"
)

var (
	servers []*gophirc.IRC
	log     = logrus.New()

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
				log.Warnln("Forcefully exiting")
				os.Exit(1)
			}
			fmt.Println()
			log.WithField("reason", "SIGINT").Infoln("Quitting")

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
		log.Fatalln(err)
	}

	err = conf.Check()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = api_config.Parse(*apiConfigFlag)
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	for _, server := range conf.Servers {
		irc := gophirc.New(server, &wg)
		err = irc.Connect()
		if err != nil {
			log.Fatalln(err)
		}

		servers = append(servers, irc)

		bot.AddBasicCallbacks(irc)
		bot.AddCTCPCallbacks(irc)

		go irc.Loop()
	}
	wg.Wait()
}
