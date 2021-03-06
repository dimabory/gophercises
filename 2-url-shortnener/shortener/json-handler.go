package shortener

import (
	"encoding/json"
	"log"
)

type pathUrlJson struct {
	Path string
	URL string
}

func NewJsonUrlMapper(filename string) (func(string) (string, bool), error) {
	content, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	var parsedUrls []pathUrlJson
	err = json.Unmarshal(content, &parsedUrls)
	if err != nil {
		log.Fatalf("Cannot parse json: %s", content)
		return nil, err
	}

	mapping := make(map[string]string)

	for _, m := range parsedUrls {
		mapping[m.Path] = m.URL
	}

	return NewBaseUrlMapper(mapping), nil
}
