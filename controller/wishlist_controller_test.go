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

			if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
				t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
			}

			if res.Header["Location"][0] != "/wishlist" {
				t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", "/wishlist", res.Header["Location"][0])
			}

			var gotWishlistId int
			var gotName string
			err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
			if err != nil {
				t.Errorf("Error querying database for inserted item: %v", err)
			}

			if gotWishlistId != 1 {
				t.Errorf("Wrong wishlist_id. Expected: %d, Got %d", 1, gotWishlistId)
			}
			if gotName != testcase.testName {
				t.Errorf("Wrong item_name. Expected: %s, Got %s", testcase.testName, gotName)
			}

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

	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
	}

	if res.Header["Location"][0] != "/wishlist/add_item" {
		t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", "/wishlist/add_item", res.Header["Location"][0])
	}

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

	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
	}

	if res.Header["Location"][0] != "/wishlist/add_item" {
		t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", "/wishlist/add_item", res.Header["Location"][0])
	}

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

	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
	}

	if res.Header["Location"][0] != "/wishlist" {
		t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", "/wishlist", res.Header["Location"][0])
	}

	req = createFormRequest("/wishlist/add_item", map[string]string{"itemName": "item_two"})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
	}

	if res.Header["Location"][0] != "/wishlist" {
		t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", "/wishlist", res.Header["Location"][0])
	}

	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items").Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	if found != 2 {
		t.Errorf("Wrong number of items in database. Expected %d item(s), Got %d", 2, found)
	}

	var gotWishlistId1 int
	var gotName1 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId1, &gotName1)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	if gotWishlistId1 != 1 {
		t.Errorf("Wrong wishlist_id on item one. Expected: %d, Got %d", 1, gotWishlistId1)
	}
	if gotName1 != "item_one" {
		t.Errorf("Wrong item_name. Expected: %s, Got %s", "item_one", gotName1)
	}

	var gotWishlistId2 int
	var gotName2 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=2").Scan(&gotWishlistId2, &gotName2)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	if gotWishlistId2 != 1 {
		t.Errorf("Wrong wishlist_id on item two. Expected: %d, Got %d", 1, gotWishlistId1)
	}
	if gotName2 != "item_two" {
		t.Errorf("Wrong item_name. Expected: %s, Got %s", "item_two", gotName1)
	}
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

	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
	}

	if res.Header["Location"][0] != "/wishlist" {
		t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", "/wishlist", res.Header["Location"][0])
	}

	req = createFormRequest("/wishlist/add_item", map[string]string{"itemName": "testItem"})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
	}

	if res.Header["Location"][0] != "/wishlist" {
		t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", "/wishlist", res.Header["Location"][0])
	}

	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items").Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	if found != 2 {
		t.Errorf("Wrong number of items in database. Expected %d item(s), Got %d", 2, found)
	}

	var gotWishlistId1 int
	var gotName1 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId1, &gotName1)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	if gotWishlistId1 != 1 {
		t.Errorf("Wrong wishlist_id on item one. Expected: %d, Got %d", 1, gotWishlistId1)
	}
	if gotName1 != "testItem" {
		t.Errorf("Wrong item_name. Expected: %s, Got %s", "item_one", gotName1)
	}

	var gotWishlistId2 int
	var gotName2 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=2").Scan(&gotWishlistId2, &gotName2)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	if gotWishlistId2 != 1 {
		t.Errorf("Wrong wishlist_id on item two. Expected: %d, Got %d", 1, gotWishlistId1)
	}
	if gotName2 != "testItem" {
		t.Errorf("Wrong item_name. Expected: %s, Got %s", "item_two", gotName1)
	}
}
