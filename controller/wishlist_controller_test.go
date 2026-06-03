package controller_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/moltenwolfcub/giftsWebApp/controller"
	"github.com/moltenwolfcub/giftsWebApp/database"
)

func genDataBase(t *testing.T) *sql.DB {
	db, err := database.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connectiong to database: %v", err)
	}
	err = database.MigrateDB(db)
	if err != nil {
		t.Errorf("Error migrating database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func createFormRequest(target string, contents map[string]string) *http.Request {
	form := url.Values{}
	for k, v := range contents {
		form.Set(k, v)
	}

	body := strings.NewReader(form.Encode())
	req := httptest.NewRequest("POST", target, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

func assertRedirect(t *testing.T, res *http.Response, expected string) {
	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
	}

	if res.Header["Location"][0] != expected {
		t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", expected, res.Header["Location"][0])
	}

}

func assert[K comparable](t *testing.T, got, want K, msg string) {
	if got != want {
		t.Errorf("%s. Expected: %v, Got %v", msg, want, got)
	}

}

func TestAddItem(t *testing.T) {
	var tests = []struct {
		testName string
		itemName string
	}{
		{
			"Normal",
			"TestItem",
		},
		{
			"Space",
			"Test Item",
		},
		{
			"Numeric",
			"3141592653589793",
		},
		{
			"Unicode",
			"!#$%& ̄f ̅ ̆ ̇ ̈ ̉ ̊ ̋ ̌ ̍ ̎ ̏ ̐ ̑ ̒ ̓ ̔ ̕ ζδψφΔΓΨΣΩ",
		},
		{
			"SQL Injection",
			"foo);DROP TABLE wishlist_items;",
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.testName, func(t *testing.T) {

			db := genDataBase(t)
			db.Exec("INSERT INTO wishlists DEFAULT VALUES")

			// cfg := config.New()
			router := controller.BuildRouter(db)

			req := createFormRequest("/wishlist/add_item", map[string]string{"itemName": testcase.testName})
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			res := w.Result()
			t.Cleanup(func() {
				res.Body.Close()
			})

			assertRedirect(t, res, "/wishlist")

			var gotWishlistId int
			var gotName string
			err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
			if err != nil {
				t.Errorf("Error querying database for inserted item: %v", err)
			}

			assert(t, gotWishlistId, 1, "Wrong wishlist_id")
			assert(t, gotName, testcase.testName, "Wrong item_name")
		})
	}
}

func TestAddItemEmpty(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/add_item", map[string]string{"itemName": ""})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist/add_item")

	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items").Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	if found > 0 {
		t.Errorf("Empty item put in database. Expected %d item(s), Got %d", 0, found)
	}
}

func TestAddItemNoBody(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := httptest.NewRequest("POST", "/wishlist/add_item", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist/add_item")

	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items").Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	if found > 0 {
		t.Errorf("Empty body POST put item in database. Expected %d item(s), Got %d", 0, found)
	}
}

func TestAddItemMultiple(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/add_item", map[string]string{"itemName": "item_one"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	req = createFormRequest("/wishlist/add_item", map[string]string{"itemName": "item_two"})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items").Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	assert(t, found, 2, "Wrong number of items in database")

	var gotWishlistId1 int
	var gotName1 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId1, &gotName1)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	assert(t, gotWishlistId1, 1, "Wrong wishlist_id")
	assert(t, gotName1, "item_one", "Wrong item_name")

	var gotWishlistId2 int
	var gotName2 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=2").Scan(&gotWishlistId2, &gotName2)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	assert(t, gotWishlistId2, 1, "Wrong wishlist_id")
	assert(t, gotName2, "item_two", "Wrong item_name")
}

func TestAddItemMultipleIdentical(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/add_item", map[string]string{"itemName": "testItem"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	req = createFormRequest("/wishlist/add_item", map[string]string{"itemName": "testItem"})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items").Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	assert(t, found, 2, "Wrong number of items in database")

	var gotWishlistId1 int
	var gotName1 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId1, &gotName1)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	assert(t, gotWishlistId1, 1, "Wrong wishlist_id")
	assert(t, gotName1, "testItem", "Wrong item_name")

	var gotWishlistId2 int
	var gotName2 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=2").Scan(&gotWishlistId2, &gotName2)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	assert(t, gotWishlistId2, 1, "Wrong wishlist_id")
	assert(t, gotName2, "testItem", "Wrong item_name")
}

func TestEditItem(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/edit_item/1", map[string]string{"itemName": "changed"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items").Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	assert(t, found, 1, "Wrong number of items in wishlist_items")

	var gotWishlistId int
	var gotName string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "changed", "Wrong item_name")
}
