CREATE TABLE wishlists (
	id INT AUTO_INCREMENT PRIMARY KEY
);
CREATE TABLE wishlist_items (
	id INT AUTO_INCREMENT PRIMARY KEY,
	wishlist_id INT NOT NULL REFERENCES wishlists (id),
	item_name TEXT NOT NULL
);
