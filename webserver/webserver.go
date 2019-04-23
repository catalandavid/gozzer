package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

const serverport string = "3030"

var htmlStatic = `<!DOCTYPE html>
<html>
<head>
    <title>Browsing WebServer</title>
</head>
<body>
	<h3>Dynamic Browsing Generator</h3>
	<ul>
		{{ range $link := . }}
			<li><a href="{{ $link }}">{{ $link }}</a></li>
		{{ end }}
	</ul>
</body>
</html>
`

func reqhandler(w http.ResponseWriter, r *http.Request) {
	urlpath := strings.TrimRight(r.URL.Path, "/")
	pathlastitem := path.Base(urlpath)
	links := []string{}
	// slashes := strings.Count(urlpath, "/")
	// nextpathitemlen := len(pathlastitem) - 1

	for i := len(pathlastitem); i > 1; i-- {
		word := []byte(pathlastitem[0 : i-1])
		url := fmt.Sprintf("%s/%s", urlpath, word)
		links = append(links, url)
	}

	tmpl, err := template.New("html").Parse(htmlStatic)

	// Error checking
	if err != nil {
		log.Panic(err)
	}

	err = tmpl.Execute(w, links)

	// Error checking
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	http.HandleFunc("/", reqhandler)

	log.Fatal(http.ListenAndServe(":"+serverport, nil))
}
