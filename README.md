# gophircbot
A simple IRC bot written from scratch*, in Go.
Uses [gophirc](https://github.com/vlad-s/gophirc).

_* Initially the bot was developed with full functionality, extracting the core into a framework was a todo._

## Description
The bot is currently a work in progress.

## Framework managed events 
* See [gophirc](https://github.com/vlad-s/gophirc)

## Bot managed events
* Searches the messages (`PRIVMSG`) for links and sends back the title or content-type & content-length (if specified in the headers)
```
       <you> https://i.imgur.com/HkN2lB4.png
<gophircbot> [URL] content-type image/png, content-length 33.4 KB
       <you> https://github.com/vlad-s/gophircbot
<gophircbot> [URL] GitHub - vlad-s/gophircbot: gophircbot, an IRC bot written from scratch, in Go
```
* Shrugs with you
```
   <someone> climate change is a hoax
       <you> shrug
<gophircbot> ¯\_(ツ)_/¯
```
* Responds to standard CTCP `VERSION`, `TIME`, `PING`, and custom `RAW` and `QUIT` - the custom CTCP events can only be invoked by a bot admin


## To do
- [ ] **Add Go documentation**
- [ ] Add back user commands
  - [x] Say, yell
  - [x] Join, part, invite, kick, ban/unban
  - [ ] Weather
  - [x] giphy API
- [x] **Extract the core functionality into a framework** - see [gophirc](https://github.com/vlad-s/gophirc)
- [ ] Add command line flags/params
  - [x] For config
  - [x] For API config
- [x] Add methods for event callbacks
- [ ] Connect the bot to a database, store info there
- [ ] Load auto replies from an external source
- [ ] Many *(?)* more