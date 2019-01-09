package shortener

import (
	"gopkg.in/yaml.v2"
	"log"
)

type pathUrlYml struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}

func NewYamlUrlMapper(filename string) (func(string) (string, bool), error) {
	content, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	var parsedUrls []pathUrlYml
	err = yaml.Unmarshal(content, &parsedUrls)
	if err != nil {
		log.Fatalf("Cannot parse yaml: %s", content)
		return nil, err
	}

	mapping := make(map[string]string)

	for _, m := range parsedUrls {
		mapping[m.Path] = m.URL
	}

	return NewBaseUrlMapper(mapping), nil
}
