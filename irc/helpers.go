package irc

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/c2h5oh/datasize"
	"github.com/pkg/errors"
	"github.com/vlad-s/gophircbot/config"
)

type User struct {
	Nick string
	User string
	Host string
}

func (u User) String() string {
	return fmt.Sprintf("%s!%s@%s", u.Nick, u.User, u.Host)
}

func (u User) IsAdmin() bool {
	for _, v := range config.Get().Admins {
		if v == u.Nick {
			return true
		}
	}
	return false
}

func ParseUser(u string) (*User, bool) {
	if u[0] == ':' {
		u = u[1:]
	}
	nb := strings.Index(u, "!")
	ub := strings.Index(u, "@")
	if nb == -1 || ub == -1 {
		return nil, false
	}
	return &User{
		Nick: u[:nb],
		User: u[nb+1 : ub],
		Host: u[ub+1:],
	}, true
}

func IsValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	return true
}

func GetTitle(u string) (title string, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", errors.Wrap(err, "Can't create new request")
	}

	req.Header.Set("User-Agent", VERSION)
	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Can't make request")
	}

	content_type, ok := res.Header["Content-Type"]
	if ok && !strings.Contains(content_type[0], "text/html") {
		title = fmt.Sprintf("content-type %s", content_type[0])

		content_length, ok := res.Header["Content-Length"]
		fmt.Printf("Content length %v; %+v\n", ok, content_length)

		if ok && content_length[0] != "" {
			parsed_size, err := strconv.ParseInt(content_length[0], 10, 64)
			if err != nil {
				return title, errors.Wrap(err, "Can't parse content length")
			}
			size := datasize.ByteSize(parsed_size)
			title += fmt.Sprintf(", content-length %s", size.HumanReadable())
		}

		return
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", errors.Wrap(err, "Can't make document from response")
	}

	title = fmt.Sprintf("%q", doc.Find("title").Text())
	title = strings.Trim(title, "\"")

	if len(title) > 150 {
		title = title[:150] + " ..."
	}

	if title == "" {
		title = "[no title]"
	}

	return
}

func IsCTCP(s string) bool {
	if s[0] == '\001' && s[len(s)-1:][0] == '\001' {
		return true
	}
	return false
}
