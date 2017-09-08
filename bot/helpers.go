package bot

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
)

const (
	VERSION = "gophircbot - See https://github.com/vlad-s/gophircbot"
)

func IsValidURL(u string) bool {
	if len(u) < 7 { // len("http://")
		return false
	}
	if u[:4] != "http" {
		return false
	}
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

	title = doc.Find("title").Text()
	title = strings.Replace(title, "\r", "", -1)
	title = strings.Replace(title, "\n", "", -1)

	if len(title) > 150 {
		title = title[:150] + " ..."
	}

	if title == "" {
		title = "[no title]"
	}

	return
}
