package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/moltenwolfcub/giftsWebApp/models"
)

type wishlistController struct {
	db *sql.DB
}

func (c *wishlistController) index(w http.ResponseWriter, r *http.Request) {
	// userID := extractUserID(r.Context())

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

	if name != "" {
		_, err := c.db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, ?);", name)

		if err != nil {
			log.Printf("Error adding wishlist item to database: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/wishlist", http.StatusFound)
	} else {
		http.Redirect(w, r, "/wishlist/add_item", http.StatusFound)
	}

}

func (c *wishlistController) editItem(w http.ResponseWriter, r *http.Request) {
	editID := r.PathValue("id")

	item, err, code := models.LoadWishlistItemStringID(c.db, editID)
	if err != nil {
		log.Print("Error loading wishlist item to edit:", err)
		http.Error(w, err.Error(), code)
		return
	}

	serveTemplate(w, "wishlist_edit.html", item)
}

func (c *wishlistController) editItemSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("itemName")
	editID := r.PathValue("id")

	item, err, code := models.LoadWishlistItemStringID(c.db, editID)
	if err != nil {
		log.Print("Error loading wishlist item to edit:", err)
		http.Error(w, err.Error(), code)
		return
	}

	if name != "" {
		item.Name = name
		err := item.Save(c.db)
		if err != nil {
			log.Printf("Error editing wishlist item in database: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/wishlist", http.StatusFound)
	} else {
		http.Redirect(w, r, "/wishlist/edit_item/"+editID, http.StatusFound)
	}
}

func (c *wishlistController) deleteItem(w http.ResponseWriter, r *http.Request) {
	deleteID := r.PathValue("id")

	var found int
	err := c.db.QueryRow("SELECT COUNT(*) FROM wishlist_items WHERE id=?", deleteID).Scan(&found)
	if err != nil {
		log.Print("Error counting items in wishlist_items:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if found != 1 {
		http.Error(w, "Tried to delete an item that doesn't exist", http.StatusNotFound)
		return
	}

	_, err = c.db.Exec("DELETE FROM wishlist_items WHERE id=?", deleteID)
	if err != nil {
		log.Printf("Error eleting wishlist item from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/wishlist", http.StatusFound)
}
