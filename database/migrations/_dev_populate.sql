INSERT INTO users (username, password, salt) VALUES ("dev", "tmpPass", "tmpSalt");

INSERT INTO wishlists (user_id) VALUES (1);
INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, "game");
INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, "toy");
