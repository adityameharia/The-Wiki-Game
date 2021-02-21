package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	start := os.Args[0]
	end := os.Args[1]
	limit := make(chan int, 250)

	pageFound := make(chan []string)

	var path []string

	set := make(map[string]bool)

	findPage(start, end, path, pageFound, limit, set)
	found := <-pageFound
	found = append([]string{start}, found...)
	for _, urlAndNextStep := range found {
		fmt.Println(urlAndNextStep)
	}
}

func getURL(url string) (pageData string) {
	res, err := http.Get(url)

	if err != nil {
		//fmt.Println(err)
		return
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return (string(bytes))
}

func urls(html string) (wikiURLs []string) {
	wikiURLs = make([]string, 0)
	wikipediaURLPrefix := "https://en.wikipedia.org"

	untrimmedAnchors := strings.Split(html, "<a href=\"")[1:]
	for _, untrimmedAnchor := range untrimmedAnchors {
		url := untrimmedAnchor[:strings.Index(untrimmedAnchor, "\"")]

		if strings.HasPrefix(url, "/wiki") && !strings.Contains(url, ":") {
			fullURL := wikipediaURLPrefix + url
			name := untrimmedAnchor[strings.Index(untrimmedAnchor, ">")+1 : strings.Index(untrimmedAnchor, "<")]
			if !strings.Contains(name, "<") || name == "Read" {
				wikiURLs = append(wikiURLs, fullURL)
			}
		}
	}

	return wikiURLs
}

func findPage(start string, end string, path []string, pageFound chan []string, limit chan int, set map[string]bool) {
	limit <- 1
	pageHTML := getURL(start)
	urls := urls(pageHTML)

	temp := make([]string, len(path))
	copy(temp, path)
	time.Sleep(2 * time.Second)
	for _, url := range urls {
		if set[url] {
			return
		}
		//fmt.Println(i)
		select {
		case <-pageFound:
			return
		default:
			temp = append(path, url)
			if url == end {
				pageFound <- temp
				close(pageFound)
				return
			}
			if len(temp) < 6 {
				go findPage(url, end, temp, pageFound, limit, set)
			}

		}
	}
	<-limit
	//fmt.Println(ch)
}
