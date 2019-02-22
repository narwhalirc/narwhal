package tusk

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/lrstanley/girc"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// NarwhalUrlParser is our url parser
var NarwhalUrlParser NarwhalUrlParserPlugin

// NewClient will create a new request-specific client, with our defined user agent, for the purposes of page fetching.
// If successful, it will return both the client and the request for use
func (parser *NarwhalUrlParserPlugin) NewClient(u url.URL) (http.Client, http.Request) {
	var requestHeaders = make(http.Header)

	client := http.Client{
		Timeout: time.Second * 15, // 15 seconds
	}

	requestHeaders.Set("User-Agent", "Narwhal Bot 0.1-alpha")

	if strings.HasSuffix(u.Host, "reddit.com") { // If this is reddit.com
		if u.Host != "old.reddit.com" { // If this isn't old.reddit.com, which is far more parser friendly
			u.Host = "old.reddit.com"
		}
	}

	request := http.Request{
		Header: requestHeaders,
		Method: "GET",
		URL:    &u,
	}

	return client, request
}

func (parser *NarwhalUrlParserPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	links := []NarwhalLink{}
	urls := []url.URL{}
	splitMessage := strings.Split(m.Message, " ") // Split on space, URLs should not contain whitespace

	for _, subMessage := range splitMessage {
		if url, parseErr := url.Parse(subMessage); parseErr == nil { // If we successfully parsed this URL
			urls = append(urls, *url) // Add this url to our urls
		}
	}

	if len(urls) > 0 { // If we have URLs
		for _, url := range urls {
			client, request := parser.NewClient(url)

			if response, getErr := client.Do(&request); getErr == nil { // If we successfully got the page
				if response.StatusCode == 200 { // Page exists
					if strings.HasPrefix(response.Header.Get("content-type"), "text/html") { // If this is likely to be an HTML page
						pageContent, readErr := ioutil.ReadAll(response.Body) // Read the body
						response.Body.Close()

						if readErr == nil { // If we successfully read page content
							doc, newDocErr := goquery.NewDocumentFromReader(bytes.NewReader(pageContent))

							if newDocErr == nil { // If we successfully got the document
								var isReddit bool
								var isYoutube bool
								var votes NarwhalRedditVotes

								title := doc.Find("title").Text() // Get the title of the page

								if strings.HasSuffix(url.Host, "reddit.com") { // If this is Reddit
									isReddit = true
									votes = NarwhalRedditVotes{
										Dislikes: doc.Find(".unvoted > .dislikes").Text(),
										Likes:    doc.Find(".unvoted > .likes").Text(),
										Score:    doc.Find(".unvoted > .unvoted").Text(),
									}
								} else if strings.HasSuffix(url.Host, "youtube.com") || strings.HasSuffix(url.Host, "youtu.be") { // If this is Youtube
									isYoutube = true
									title = strings.TrimSuffix(title, "- YouTube") // Strip - Youtube from title
								}

								links = append(links, NarwhalLink{
									IsReddit:  isReddit,
									IsYoutube: isYoutube,
									Link:      url,
									Title:     strings.TrimSpace(title),
									Votes:     votes,
								})
							}
						}
					}
				}
			}
		}
	}

	if len(links) > 0 { // If we found links
		for _, link := range links { // For each link
			if link.IsReddit { // If this is Reddit
				parser.ParseReddit(c, e, link) // Hand off to ParseReddit
			} else if link.IsYoutube { // If this is Youtube
				parser.ParseYoutube(c, e, link) // Hand off to ParseYoutube
			} else { // Some other link
				title := fmt.Sprintf("[ %s ]", link.Title)
				c.Cmd.Reply(e, title)
			}
		}
	}
}

// ParseReddit will parse the Reddit URL and metadata, outputting into the event the information
func (parser *NarwhalUrlParserPlugin) ParseReddit(c *girc.Client, e girc.Event, l NarwhalLink) {
	var title string

	if l.Votes.Dislikes != "" && l.Votes.Likes != "" && l.Votes.Score != "" {
		var convertScoreErr error
		var downvotes int
		var upvotes int

		downvotes, convertScoreErr = strconv.Atoi(l.Votes.Dislikes)

		if convertScoreErr == nil { // No error converting downvotes
			upvotes, convertScoreErr = strconv.Atoi(l.Votes.Likes)
		}

		if convertScoreErr == nil { // No error converting downvotes and upvotes
			percentage := int((float64(downvotes) / float64(upvotes)) * 100)

			if percentage == 0 { // 100% upvote
				percentage = 100
			}

			title = fmt.Sprintf("[ %s ][Score: %s, %d%% upvotes]", l.Title, l.Votes.Score, percentage)
		} else {
			title = fmt.Sprintf("[ %s ]", l.Title)
		}
	} else {
		title = fmt.Sprintf("[ %s ]", l.Title)
	}

	c.Cmd.Reply(e, title)
}

// ParseYoutube will parse the Youtube URL and output into the event target the information
func (parser *NarwhalUrlParserPlugin) ParseYoutube(c *girc.Client, e girc.Event, l NarwhalLink) {
	host := l.Link.Hostname()
	queriesMap := l.Link.Query() // Get the queries

	var title string
	var isVideo bool
	var vidUrl string

	if strings.HasSuffix(host, "youtu.be") { // URL shortened version
		isVideo = true
		vidUrl = strings.TrimPrefix(l.Link.Path, "/") // Trim the leading / from path first
	} else { // Non-shortened
		if vidUrls, hasVQueries := queriesMap["v"]; hasVQueries {
			if len(vidUrls) > 0 { // If this has a ?v=
				isVideo = true
				vidUrl = vidUrls[0] // Only get first one
			}
		}
	}

	if isVideo { // If this is a Youtube video
		desktopYT := fmt.Sprintf("https://youtube.com/watch?v=%s", vidUrl)
		mobileYT := fmt.Sprintf("https://m.youtube.com/watch?v=%s", vidUrl)
		title = fmt.Sprintf("[ %s | Desktop: %s | Mobile: %s ]", l.Title, desktopYT, mobileYT)
	} else { // Not a Youtube video
		title = fmt.Sprintf("[ %s ]", l.Title)
	}

	c.Cmd.Reply(e, title) // Reply to target with title
}
