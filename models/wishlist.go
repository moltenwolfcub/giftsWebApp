package models

import (
	"database/sql"
	"log"
)

type WishlistItem struct {
	Id   int
	Name string
}

func LoadWishlistItem(db *sql.DB, id int) (*WishlistItem, error) {
	var name string
	err := db.QueryRow("SELECT item_name FROM wishlist_items WHERE id=?", id).Scan(&name)
	if err != nil {
		log.Printf("Error finding wishlist item[%d]: %v", id, err)
		return nil, err
	}

	item := WishlistItem{
		Id:   id,
		Name: name,
	}
	return &item, nil
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
		wl.Items = append(wl.Items, WishlistItem{Id: id, Name: name})
	}

	return &wl, nil
}
