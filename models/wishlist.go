package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type WishlistItem struct {
	Id   int
	Name string

	wishlistId int

	stale bool
}

func CreateWishlistItem(db *sql.DB, name string, wishlistId int) (*WishlistItem, error) {
	res, err := db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (0, \"__empty\");")
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	item := &WishlistItem{
		Id:   int(id),
		Name: name,

		wishlistId: wishlistId,
	}
	return item, nil
}

func LoadWishlistItemStringID(db *sql.DB, id string) (*WishlistItem, error, int) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Print("Non-integer id given to editItem handler:", err)
		return nil, err, http.StatusNotFound
	}
	return LoadWishlistItem(db, intId)
}

func LoadWishlistItem(db *sql.DB, id int) (*WishlistItem, error, int) {
	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items WHERE id=?", id).Scan(&found)
	if err != nil {
		log.Print("Error counting items in wishlist_items:", err)
		return nil, err, http.StatusInternalServerError
	}
	if found == 0 {
		return nil, fmt.Errorf("Requested wishlist item [%d] doesn't exist.", id), http.StatusNotFound
	}
	if found > 1 {
		return nil, fmt.Errorf("Multiple database entries containing wishlist item [%d]", id), http.StatusNotFound
	}

	var name string
	var wishlistId int
	err = db.QueryRow("SELECT item_name, wishlist_id FROM wishlist_items WHERE id=?", id).Scan(&name, &wishlistId)
	if err != nil {
		log.Printf("Error finding wishlist item[%d]: %v", id, err)
		return nil, err, http.StatusInternalServerError
	}

	item := WishlistItem{
		Id:         id,
		Name:       name,
		wishlistId: wishlistId,
	}
	return &item, nil, http.StatusOK
}

func (w *WishlistItem) Save(db *sql.DB) error {
	if w.stale {
		return fmt.Errorf("Tried to save an item that isn't alive")
	}

	_, err := db.Exec("UPDATE wishlist_items SET item_name=?,wishlist_id=? WHERE id=?", w.Name, w.wishlistId, w.Id)
	if err != nil {
		return err
	}
	return nil
}

func (w *WishlistItem) Delete(db *sql.DB) error {
	if w.stale {
		return fmt.Errorf("Tried to delete an item that already isn't alive")
	}

	w.stale = true
	_, err := db.Exec("DELETE FROM wishlist_items WHERE id=?", w.Id)
	if err != nil {
		return err
	}
	return nil
}

type Wishlist struct {
	Items []WishlistItem
}

func LoadWishlist(db *sql.DB) (*Wishlist, error) {
	wl := Wishlist{}
	wl.Items = []WishlistItem{}

	var listId int
	err := db.QueryRow("SELECT id FROM wishlists LIMIT 1").Scan(&listId)
	if err != nil {
		log.Printf("Error finding wishlist: %v", err)
		return nil, err
	}

	items, err := db.Query("SELECT id, item_name FROM wishlist_items WHERE wishlist_id=?", listId)
	if err != nil {
		log.Printf("Error finding wishlist items: %v", err)
		return nil, err
	}
	defer items.Close()

	for items.Next() {
		var name string
		var id int
		items.Scan(&id, &name)
		wl.Items = append(wl.Items, WishlistItem{
			Id:         id,
			Name:       name,
			wishlistId: 1, //TODO when implementing multiple wishlists change this
		})
	}

	return &wl, nil
}
