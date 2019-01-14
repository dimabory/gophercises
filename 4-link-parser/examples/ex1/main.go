package main

import (
	"fmt"
	"github.com/dimabory/gophercises/4-link-parser"
	"strings"
)

var exampleHtml = `
<html>
<body>
    <h1>Hello!</h1>
    <a href="/other-page">A link to another page</a>
    <div>1231</div>
    <a href="/sec-page">blablablab</a>
    <a href="/ssss-qqq">
		asdasd 
		<span>asdasd</span> 
		asd xz,cvmnlodsr
	</a>
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
