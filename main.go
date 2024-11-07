package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

var target_url = "https://gobyexample.com"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Regexphref(expression string) (*regexp.Regexp, error) {
	re, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	return re, nil
}

func GetHTMLResponse(target_url string) ([]byte, error) {
	response, err := http.Get(target_url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d", response.StatusCode)
	}

	return io.ReadAll(response.Body)
}

func ParesURLs(html string, target_url string, url_queue URLQueue) {
	re_a_href, err := Regexphref(`a href="([^"]*)"`)
	check(err)
	base, err := url.Parse(target_url)
	check(err)
	for _, link := range re_a_href.FindAllStringSubmatch(html, -1) {
		l := link[1]
		parsedURL, err := url.Parse(l)
		if err != nil {
			fmt.Println("err parsing url:", err)
			continue
		}
		absoluteURL := base.ResolveReference(parsedURL).String()
		url_queue.AddURL(absoluteURL)
	}
}

func Webcrawler(target_url string, url_queue URLQueue) {
	url_queue.AddURL(target_url)
	html, err := GetHTMLResponse(target_url + "/")
	check(err)
	ParesURLs(string(html), target_url, url_queue)
}

type URLQueue map[string]bool

func (q URLQueue) HasURL(link string) bool {
	_, ok := q[link]
	return ok
}
func (q URLQueue) AddURL(link string) {
	if q.HasURL(link) {
		q[link] = true
	} else {
		q[link] = false
	}
}

func main() {
	url_queue := URLQueue{}
	Webcrawler(target_url, url_queue)
	fmt.Println(url_queue)
}
