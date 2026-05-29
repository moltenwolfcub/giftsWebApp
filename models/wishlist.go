package models

import "database/sql"

type WishlistItem struct {
	Id   int
	Name string
}

type Wishlist struct {
	Items []WishlistItem
}

func LoadWishlist(db *sql.DB) (*Wishlist, error) {
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
