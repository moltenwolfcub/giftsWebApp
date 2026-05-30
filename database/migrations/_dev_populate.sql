INSERT INTO users (username, password, salt) VALUES (
	"dev",
	"nil",
	"devSalt"
);
UPDATE users SET password=X'3a4d722d5f2f1908bbc96eeeafcacebb4cc53187698f75d07de9fd144adc9681'
WHERE username="dev";

INSERT INTO wishlists (user_id) VALUES (1);
INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, "game");
INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, "toy");
