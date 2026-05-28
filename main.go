package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const templatesPath = "templates"

type Wishlist struct {
	Items []string
}

var dummyWishlist = Wishlist{[]string{"Chocolate", "Sweets", "Guitar"}}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/wishlist/", handleWishlist)
	log.Fatal(http.ListenAndServe(":8040", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home\nHello World")
}

func handleWishlist(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "wishlist.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, dummyWishlist)
}
