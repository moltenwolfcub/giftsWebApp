package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/wishlist/", handleWishlist)
	log.Fatal(http.ListenAndServe(":8040", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home\nHello World")
}

func handleWishlist(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Your Wishlist\nHello World")
}
