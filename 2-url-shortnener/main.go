package main

import (
	"flag"
	"fmt"
	"github.com/dimabory/gophercises/2-url-shortnener/shortener"
	"log"
	"net/http"
)

var (
	pathToUrls        map[string]string
	mux               *http.ServeMux
	yamlMappingPath   string
	jsonMappingPath   string
	boltDbMappingPath string
)

func init() {
	flag.StringVar(&yamlMappingPath, "yaml", "./mapping.yml", "yaml mapping file path")
	flag.StringVar(&jsonMappingPath, "json", "./mapping.json", "yaml mapping file path")
	flag.StringVar(&boltDbMappingPath, "boltdb", "./mapping.db", "yaml mapping file path")

	flag.Parse()

	pathToUrls = map[string]string{
		"/godoc":     "https://godoc.org/github.com/gophercises/urlshort",
		"/godoc.yml": "https://godoc.org/gopkg.in/yaml.v2",
	}

	mux = defaultMux()
}

func main() {
	defaultUrlMapper := shortener.NewBaseUrlMapper(pathToUrls)
	redirectHandler := shortener.NewHttpRedirectHandler(defaultUrlMapper, mux)

	ymlUrlMapper, err := shortener.NewYamlUrlMapper(yamlMappingPath)
	if err != nil {
		log.Fatalf("Can't create YAML redirect URL provider. %v", err)
	}

	redirectHandler = shortener.NewHttpRedirectHandler(ymlUrlMapper, redirectHandler)

	jsonUrlMapper, err := shortener.NewJsonUrlMapper(jsonMappingPath)
	if err != nil {
		log.Fatalf("Can't create JSON redirect URL provider. %v", err)
	}

	redirectHandler = shortener.NewHttpRedirectHandler(jsonUrlMapper, redirectHandler)

	boltdbMapper, err := shortener.NewBoltDbUrlMapper(boltDbMappingPath, "pathes")
	if err != nil {
		log.Fatalf("Can't create BoltDB redirect URL provider. %v", err)
	}

	redirectHandler = shortener.NewHttpRedirectHandler(boltdbMapper, redirectHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", redirectHandler)
}

func defaultMux() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", hello)
	return
}

func hello(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Hello, world!")
}
