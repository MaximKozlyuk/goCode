package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	page string
)

func main() {

	resp, _ := http.Get("https://ru.wikipedia.org/wiki/Заглавная_страница")
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	body_strings := string(bytes)
	page = body_strings

	//data, _ := ioutil.ReadFile("wikiPage.html")

	http.HandleFunc("/", index_handler)
	http.ListenAndServe(":8000", nil)

	fmt.Println("Done")
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, page)
}

type SiteMapIndex struct {
	Locations []string
}

type News struct {
	Titles    []string
	Keywords  []string
	Locations []string
}
