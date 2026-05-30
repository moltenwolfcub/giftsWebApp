INSERT INTO users (username, password, salt) VALUES (
	"dev",
	"9265886dd625b0332b04a6633f6131a55b7176023cf148b2fcea67946d2e68e3",
	"7bf2fd255f5e9e09ebb2d818c98f45ff"
);

INSERT INTO wishlists (user_id) VALUES (1);
INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, "game");
INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, "toy");
