package adventure

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func init() {
	defaultTemplate = template.Must(template.New("").Parse(htmlTemplate))
}

func ReadFile(filename string) (io.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return file, err
}

func JsonStory(file io.Reader) (story Story, err error) {
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}

	return
}

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.template = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func NewHandler(story Story, opts ...HandlerOption) http.Handler {
	h := handler{story, defaultTemplate, defaultPathFn}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

type HandlerOptions struct {
	*template.Template
	ParseFunc func(r *http.Request) string
}

type handler struct {
	story    Story
	template *template.Template
	pathFn   func(r *http.Request) string
}

func defaultPathFn(r *http.Request) (path string) {
	path = strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	return
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.story[path]; ok {

		if err := h.template.Execute(w, chapter); err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

var defaultTemplate *template.Template

var htmlTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>
`
