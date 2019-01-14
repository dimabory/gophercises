package main

import (
	"fmt"
	"github.com/dimabory/gophercises/4-link-parser"
	"strings"
)

var exampleHtml = `
<html>
<body>
    <a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
</body>
</html>
`

func main() {
	r := strings.NewReader(string(exampleHtml))

	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", links)
}
