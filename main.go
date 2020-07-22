// Copyright 2013 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"encoding/hex"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var templates map[string]*template.Template

func checkFor500s(err error) {
	if err != nil {
		panic(err)
	}
}

func handlerFor500s(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := templates["home"].ExecuteTemplate(w, "base", nil)
	checkFor500s(err)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "HTTP method is not POST", http.StatusNotFound)
	}

	f, _, err := r.FormFile("file")
	checkFor500s(err)
	defer f.Close()

	var encoder io.Writer
	switch strings.ToLower(r.FormValue("encoding")) {
	case "base64":
		encoder = base64.NewEncoder(base64.URLEncoding, w)
	default:
		encoder = hex.NewEncoder(w)
	}

	w.Header().Add("Content-Disposition", "attachment; filename=catbyte.hex.txt")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(encoder, f)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortPath := vars["path"]
	path := "static/" + shortPath
	fileinfo, err := os.Stat(path)
	if err != nil || fileinfo.IsDir() {
		http.Error(w, "Unable to find file", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}

func init() {
	// Templates
	templates = make(map[string]*template.Template)
	templates["home"] = template.Must(template.ParseFiles("templates/home.html", "templates/base.html"))
	templates["404"] = template.Must(template.ParseFiles("templates/404.html", "templates/base.html"))
}

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/static").Subrouter()
	s.HandleFunc("/{path:.*?}", handlerFor500s(StaticHandler)).Name("static")
	r.HandleFunc("/upload", handlerFor500s(UploadHandler)).Name("image")
	r.HandleFunc("/", handlerFor500s(HomeHandler)).Name("home")

	http.Handle("/", r)
	http.ListenAndServe(":80", nil)
}
