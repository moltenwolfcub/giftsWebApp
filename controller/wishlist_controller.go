package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/moltenwolfcub/giftsWebApp/models"
)

type wishlistController struct {
	db *sql.DB
}

func (c *wishlistController) index(w http.ResponseWriter, r *http.Request) {
	wishlist, err := models.LoadWishlist(c.db)
	if err != nil {
		log.Print("Error loading wishlist from database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serveTemplate(w, "wishlist.html", wishlist)
}

func (c *wishlistController) addItem(w http.ResponseWriter, r *http.Request) {
	serveTemplate(w, "wishlist_add.html", nil)
}

func (c *wishlistController) addItemSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("itemName")

	_, err := c.db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, ?);", name)
	if err != nil {
		log.Printf("Error adding wishlist item to database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}

func (c *wishlistController) editItem(w http.ResponseWriter, r *http.Request) {
	editID := r.PathValue("id")
	editIDint, err := strconv.Atoi(editID)
	if err != nil {
		log.Print("Non-integer id given to editItem handler:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	item, err := models.LoadWishlistItem(c.db, editIDint)
	if err != nil {
		log.Print("Error loading wishlist item to edit:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serveTemplate(w, "wishlist_edit.html", item)
}

func (c *wishlistController) editItemSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("itemName")
	editID := r.PathValue("id")

	_, err := c.db.Exec("UPDATE wishlist_items SET item_name=? WHERE id=?", name, editID)
	if err != nil {
		log.Printf("Error editing wishlist item in database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}

func (c *wishlistController) deleteItem(w http.ResponseWriter, r *http.Request) {
	deleteID := r.PathValue("id")

	_, err := c.db.Exec("DELETE FROM wishlist_items WHERE id=?", deleteID)
	if err != nil {
		log.Printf("Error eleting wishlist item from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}
