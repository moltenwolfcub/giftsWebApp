package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const templatesPath = "templates"
const databasePath = "database"

type WishlistItem struct {
	Id   int
	Name string
}

type Wishlist struct {
	Items []WishlistItem
}

func LoadWishlist() (*Wishlist, error) {
	wl := Wishlist{}
	wl.Items = []WishlistItem{}

	var listId int
	err := db.QueryRow("SELECT * FROM wishlists LIMIT 1").Scan(&listId)
	if err != nil {
		return nil, err
	}

	items, err := db.Query("SELECT * FROM wishlist_items WHERE wishlist_id=?", listId)
	if err != nil {
		return nil, err
	}
	defer items.Close()

	for items.Next() {
		var name string
		var id, nil2 int
		items.Scan(&id, &nil2, &name)
		wl.Items = append(wl.Items, WishlistItem{Id: id, Name: name})
	}

	return &wl, nil
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", filepath.Join(databasePath, "dev.db"))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Connected to database")

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/wishlist/", handleWishlist)
	http.HandleFunc("/wishlist/add_item", handleAddWishlistItem)
	http.HandleFunc("POST /wishlist/add_item/submit", handleAddWishlistItemSubmit)
	http.HandleFunc("/wishlist/edit_item", handleEditWishlistItem)
	http.HandleFunc("POST /wishlist/edit_item/submit", handleEditWishlistItemSubmit)
	http.HandleFunc("POST /wishlist/delete_item", handleDeleteWishlistItem)
	log.Fatal(http.ListenAndServe(":8040", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "index.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, nil)
}

func handleWishlist(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "wishlist.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	wishlist, err := LoadWishlist()
	if err != nil {
		log.Print("Error loading wishlist from database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, wishlist)
}

func handleAddWishlistItem(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "wishlist_add.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, nil)
}

func handleAddWishlistItemSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("itemName")

	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, ?);", name)

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}

func handleEditWishlistItem(w http.ResponseWriter, r *http.Request) {
	editID := r.FormValue("id")
	var item WishlistItem
	var _nil0 string
	err := db.QueryRow("SELECT * FROM wishlist_items WHERE id=?", editID).Scan(&item.Id, &_nil0, &item.Name)
	if err != nil {
		log.Print("Error finding item to edit:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "wishlist_edit.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, item)
}

func handleEditWishlistItemSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("itemName")
	editID := r.FormValue("editID")

	db.Exec("UPDATE wishlist_items SET item_name=? WHERE id=?", name, editID)

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}

func handleDeleteWishlistItem(w http.ResponseWriter, r *http.Request) {
	deleteID := r.FormValue("id")

	db.Exec("DELETE FROM wishlist_items WHERE id=?", deleteID)

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}
