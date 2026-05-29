package controller

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/moltenwolfcub/giftsWebApp/models"
)

const templatesPath = "templates"

var db *sql.DB

func BuildRouter(database *sql.DB) {
	db = database

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/wishlist/", handleWishlist)
	http.HandleFunc("/wishlist/add_item", handleAddWishlistItem)
	http.HandleFunc("POST /wishlist/add_item/submit", handleAddWishlistItemSubmit)
	http.HandleFunc("/wishlist/edit_item", handleEditWishlistItem)
	http.HandleFunc("POST /wishlist/edit_item/submit", handleEditWishlistItemSubmit)
	http.HandleFunc("POST /wishlist/delete_item", handleDeleteWishlistItem)

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "_base.html"),
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
		filepath.Join(templatesPath, "_base.html"),
		filepath.Join(templatesPath, "wishlist.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	wishlist, err := models.LoadWishlist(db)
	if err != nil {
		log.Print("Error loading wishlist from database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, wishlist)
}

func handleAddWishlistItem(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "_base.html"),
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
	var item models.WishlistItem
	var _nil0 string
	err := db.QueryRow("SELECT * FROM wishlist_items WHERE id=?", editID).Scan(&item.Id, &_nil0, &item.Name)
	if err != nil {
		log.Print("Error finding item to edit:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "_base.html"),
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
