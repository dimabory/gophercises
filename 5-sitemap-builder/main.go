package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/dimabory/gophercises/4-link-parser"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

var (
	url      string
	maxDepth int
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

func init() {
	flag.StringVar(&url, "url", "https://godoc.org", "a domain to start building a sitemap")
	flag.IntVar(&maxDepth, "depth", 1, "the max number of links deep to traverse")
}

func main() {
	pages := bfs(url, maxDepth)

	toXml := urlset{Xmlns: xmlns,}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "	")

	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})

		if len(q) == 0 {
			break
		}

		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}

			seen[url] = struct{}{}
			for _, link := range get(url) {
				if _, ok := seen[link]; !ok {
					nq[link] = struct{}{}
				}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(path string) []string {
	fmt.Printf("Getting %s ...\n", path)
	resp, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url2.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	return filter(hrefs(resp.Body, baseUrl.String()), withPrefix(baseUrl.String()))
}

func hrefs(body io.Reader, baseUrl string) []string {
	links, err := link.Parse(body)
	if err != nil {
		panic(err)
	}
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, baseUrl+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string

	for _, l := range links {
		if keepFn(l) && !contains(ret, l) {
			ret = append(ret, l)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(l string) bool {
		return strings.HasPrefix(l, pfx)
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))

	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}
