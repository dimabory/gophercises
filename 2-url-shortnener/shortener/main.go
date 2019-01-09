package shortener

import (
	"io/ioutil"
	"log"
	"net/http"
)

func NewBaseUrlMapper(urls map[string]string) func(string) (string, bool) {
	return func(path string) (url string, ok bool) {
		url, ok = urls[path]
		return
	}
}

func NewHttpRedirectHandler(mapper func(string) (string, bool), fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestedUrl := request.URL.Path
		if url, ok := mapper(requestedUrl); ok {
			log.Printf("Redirecting %s to %s\n", requestedUrl, url)
			http.Redirect(writer, request, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(writer, request)
	}
}

func readFile(filename string) (content []byte, err error) {
	content, err = ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Cannot read file: %s", filename)
		return nil, err
	}

	return
}
