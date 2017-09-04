# gophircbot
A simple IRC bot written from scratch, in Go.

## Description
The bot is currently a work in progress; the goal is to write an easy to use, extendable, event based IRC framework for regular clients or bots.

Currently, the bot loads its config from a `config.json` file, the config's path being relative to the binary.

## Self managed events
* Auto pongs the pings
* Registers automatically on first `NOTICE *`
* Identifies automatically on `RPL_WELCOME` (event 001)
* Automatically joins the received invites & sends a greeting to the channel
* Logs if the bot gets kicked from a channel

## Features
* Events are parsed & commands are sent through goroutines
* Gets links' titles or content-type & content-length if available
* Parses CTCP events: `VERSION`, `TIME`, `PING`, and implements custom `RAW` and `QUIT`; the custom CTCP events can only be invoked by an admin user
* Parses built-in commands: `say`, `yell`, `join`, `part`; the commands **must be** preceded by a comma (`,`). The `join` and `part` commands can only be invoked by an admin user
* State logging - logs on connection, registering, identifying.
* Graceful exit through a `quit chan struct{}`, handled either by a `SIGINT` (Ctrl-C) or a `CTCP QUIT`
* Parses a user from an IRC formatted `nick!user@host` to a `User{}`
* Config implements a basic checking on values


* Many *(?)* more

## To do
* **Extract the core functionality into a framework**
* Add defaults
* Add more commands: `MODE`, `KICK`, etc.
* Add regex matching for nicknames, channels
* Add command line flags/params for config & probably other things
* Add methods for event callbacks
* Connect the bot to a database, store info there
* Load auto replies from an external source


* Many *(?)* more