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

			formContents := url.Values{}
			formContents.Set("itemName", testcase.testName)

			body := strings.NewReader(formContents.Encode())
			req := httptest.NewRequest("POST", "/wishlist/add_item", body)
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

	formContents := url.Values{}
	formContents.Set("itemName", "")

	body := strings.NewReader(formContents.Encode())
	req := httptest.NewRequest("POST", "/wishlist/add_item", body)
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

	var found int
	err := db.QueryRow("SELECT COUNT(*) FROM wishlist_items").Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	if found > 0 {
		t.Errorf("Empty item put in database. Expected %d item(s), Got %d", 0, found)
	}
}
