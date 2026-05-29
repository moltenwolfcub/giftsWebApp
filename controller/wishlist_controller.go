package controller

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/moltenwolfcub/giftsWebApp/models"
)

type wishlistController struct {
	db *sql.DB
}

func (c *wishlistController) index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "_base.html"),
		filepath.Join(templatesPath, "wishlist.html"),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	wishlist, err := models.LoadWishlist(c.db)
	if err != nil {
		log.Print("Error loading wishlist from database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, wishlist)
}

func (c *wishlistController) addItem(w http.ResponseWriter, r *http.Request) {
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

func (c *wishlistController) addItemSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("itemName")

	c.db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, ?);", name)

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}

func (c *wishlistController) editItem(w http.ResponseWriter, r *http.Request) {
	editID := r.FormValue("id")
	var item models.WishlistItem
	var _nil0 string
	err := c.db.QueryRow("SELECT * FROM wishlist_items WHERE id=?", editID).Scan(&item.Id, &_nil0, &item.Name)
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

func (c *wishlistController) editItemSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("itemName")
	editID := r.FormValue("editID")

	c.db.Exec("UPDATE wishlist_items SET item_name=? WHERE id=?", name, editID)

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}

func (c *wishlistController) deleteItem(w http.ResponseWriter, r *http.Request) {
	deleteID := r.FormValue("id")

	c.db.Exec("DELETE FROM wishlist_items WHERE id=?", deleteID)

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}
