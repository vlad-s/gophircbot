package bot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/c2h5oh/datasize"
	"github.com/pkg/errors"
	"github.com/vlad-s/gophircbot/apiconfig"
	"github.com/vlad-s/gophircbot/db"
)

const (
	// VERSION stores the bot version
	VERSION = "gophircbot - See https://github.com/vlad-s/gophircbot"
)

// IsValidURL returns whether a string is a valid URL.
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

// GetTitle returns the title of a web page, or an error in case it fails.
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

	contentTypes, ok := res.Header["Content-Type"]
	if ok && !strings.Contains(contentTypes[0], "text/html") {
		title = fmt.Sprintf("content-type %s", contentTypes[0])

		contentLengths, ok := res.Header["Content-Length"]

		if ok && contentLengths[0] != "" {
			parsedSize, err := strconv.ParseInt(contentLengths[0], 10, 64)
			if err != nil {
				return title, errors.Wrap(err, "Can't parse content length")
			}
			size := datasize.ByteSize(parsedSize)
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

// GetGif returns the giphy URL of a gif, searched by "query", or an error in case it fails.
func GetGif(query string) (reply string, err error) {
	reply = "[giphy] "

	apiConfig := apiconfig.Get()
	giphyURL := "https://api.giphy.com/v1/gifs/search?q=%s&api_key=%s&limit=%d"

	queryURL := fmt.Sprintf(giphyURL, url.QueryEscape(query), apiConfig.Giphy.APIKey, apiConfig.Giphy.Limit)
	res, err := http.Get(queryURL)
	if err != nil {
		return "", errors.Wrap(err, "Error requesting the giphy URL")
	}
	defer res.Body.Close()

	var gifs apiconfig.GiphyResponse
	err = json.NewDecoder(res.Body).Decode(&gifs)
	if err != nil {
		return "", errors.Wrap(err, "Error decoding the JSON response")
	}

	if gifs.Pagination.Total == 0 || gifs.Pagination.Count == 0 || len(gifs.Data) == 0 {
		reply += "no GIFs found :("
		return
	}

	gif := gifs.Data[rand.Intn(gifs.Pagination.Count)]

	if gif.Rating == "r" {
		reply += "NSFW "
	}

	reply += gif.ShortURL
	return
}

// IsIgnored returns whether the supplied nick is ignored or not.
// Depends only on the database, the config array is verified by the framework.
func IsIgnored(nick string) bool {
	return !db.Get().Where("nick = ?", nick).First(&db.IgnoredUser{}).RecordNotFound()
}
