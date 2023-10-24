package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

const (
	addr = "0.0.0.0"
	port = "8080"
)

var (
	//go:embed templates
	lfs embed.FS
)

type Input struct {
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodPost:
			execute(w, r)
		default:
			form(w, r)
		}
	})
	log.Println("server is running...")
	return http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), nil)
}

func form(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFS(lfs, "templates/*")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "parse template error: %s", err.Error())
		return
	}
	err = tpl.ExecuteTemplate(w,
		"form.tmpl",
		struct{ Port string }{Port: port})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "execute template error: %s", err.Error())
		return
	}
}

func execute(w http.ResponseWriter, r *http.Request) {
	var (
		input Input
		err   error
	)

	if err = r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "decode form error: %s", err.Error())
		return
	}

	input.Template = r.PostForm.Get("template")
	data := r.PostForm.Get("data")

	if err = json.Unmarshal(([]byte)(data), &input.Data); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "decode json error: %s", err.Error())
		return
	}
	tpl, err := template.New("tpl").Parse(input.Template)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "parse template error: %s", err.Error())
		return
	}
	if err = tpl.Execute(w, input.Data); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "execute template error: %s", err.Error())
		return
	}
}
