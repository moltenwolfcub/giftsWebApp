package controller_test

import (
	"database/sql"
	"fmt"
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

func assertCount(t *testing.T, db *sql.DB, table string, expected int) {
	var found int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	assert(t, found, expected, "Wrong number of items in "+table)
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

			assertCount(t, db, "wishlist_items", 1)

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

	assertCount(t, db, "wishlist_items", 0)
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

	assertCount(t, db, "wishlist_items", 0)
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

	assertCount(t, db, "wishlist_items", 2)

	var gotWishlistId1 int
	var gotName1 string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId1, &gotName1)
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

	assertCount(t, db, "wishlist_items", 2)

	var gotWishlistId1 int
	var gotName1 string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId1, &gotName1)
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
	var tests = []struct {
		testName   string
		editedName string
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
			db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

			// cfg := config.New()
			router := controller.BuildRouter(db)

			req := createFormRequest("/wishlist/edit_item/1", map[string]string{"itemName": testcase.editedName})
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			res := w.Result()
			t.Cleanup(func() {
				res.Body.Close()
			})

			assertRedirect(t, res, "/wishlist")

			assertCount(t, db, "wishlist_items", 1)

			var gotWishlistId int
			var gotName string
			err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
			if err != nil {
				t.Errorf("Error querying database for changed item: %v", err)
			}

			assert(t, gotWishlistId, 1, "Wrong wishlist_id")
			assert(t, gotName, testcase.editedName, "Wrong item_name")
		})
	}
}

func TestEditItemEmpty(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/edit_item/1", map[string]string{"itemName": ""})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist/edit_item/1")

	assertCount(t, db, "wishlist_items", 1)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "initialItem", "Wrong item_name")
}

func TestEditSameItemTwice(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	// FIRST CHANGE
	req := createFormRequest("/wishlist/edit_item/1", map[string]string{"itemName": "updatedVersion"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	assertCount(t, db, "wishlist_items", 1)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "updatedVersion", "Wrong item_name")

	// SECOND CHANGE
	req = createFormRequest("/wishlist/edit_item/1", map[string]string{"itemName": "finalVersion"})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	assertCount(t, db, "wishlist_items", 1)

	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "finalVersion", "Wrong item_name")
}

func TestEditItemWithOthers(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"secondaryItem\");")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"thirdItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/edit_item/1", map[string]string{"itemName": "editedItem"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	assertCount(t, db, "wishlist_items", 3)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "editedItem", "Wrong item_name")

	var gotWishlistId2 int
	var gotName2 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=2").Scan(&gotWishlistId2, &gotName2)
	if err != nil {
		t.Errorf("Error querying database for unchanged item: %v", err)
	}

	assert(t, gotWishlistId2, 1, "Wrong wishlist_id")
	assert(t, gotName2, "secondaryItem", "Wrong item_name")

	var gotWishlistId3 int
	var gotName3 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=3").Scan(&gotWishlistId3, &gotName3)
	if err != nil {
		t.Errorf("Error querying database for unchanged item: %v", err)
	}

	assert(t, gotWishlistId3, 1, "Wrong wishlist_id")
	assert(t, gotName3, "thirdItem", "Wrong item_name")
}

func TestEditItemMultiple(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"secondaryItem\");")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"thirdItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/edit_item/1", map[string]string{"itemName": "editedItem"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	req = createFormRequest("/wishlist/edit_item/2", map[string]string{"itemName": "otherItemEdited"})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	assertCount(t, db, "wishlist_items", 3)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "editedItem", "Wrong item_name")

	var gotWishlistId2 int
	var gotName2 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=2").Scan(&gotWishlistId2, &gotName2)
	if err != nil {
		t.Errorf("Error querying database for unchanged item: %v", err)
	}

	assert(t, gotWishlistId2, 1, "Wrong wishlist_id")
	assert(t, gotName2, "otherItemEdited", "Wrong item_name")

	var gotWishlistId3 int
	var gotName3 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=3").Scan(&gotWishlistId3, &gotName3)
	if err != nil {
		t.Errorf("Error querying database for unchanged item: %v", err)
	}

	assert(t, gotWishlistId3, 1, "Wrong wishlist_id")
	assert(t, gotName3, "thirdItem", "Wrong item_name")
}

func TestEditItemMultipleIdentical(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"anItem\");")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"differentItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/edit_item/2", map[string]string{"itemName": "anItem"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	assertCount(t, db, "wishlist_items", 2)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "anItem", "Wrong item_name")

	var gotWishlistId2 int
	var gotName2 string
	err = db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=2").Scan(&gotWishlistId2, &gotName2)
	if err != nil {
		t.Errorf("Error querying database for unchanged item: %v", err)
	}

	assert(t, gotWishlistId2, 1, "Wrong wishlist_id")
	assert(t, gotName2, "anItem", "Wrong item_name")
}

func TestEditItemInvalidID(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/edit_item/5", map[string]string{"itemName": "invalid"})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusNotFound, res.StatusCode)
	}

	assertCount(t, db, "wishlist_items", 1)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "initialItem", "Wrong item_name")
}

func TestEditItemNoBody(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := httptest.NewRequest("POST", "/wishlist/edit_item/1", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist/edit_item/1")

	assertCount(t, db, "wishlist_items", 1)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "initialItem", "Wrong item_name")
}

func TestDeleteItem(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/delete_item/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	assertCount(t, db, "wishlist_items", 0)
}

func TestDeleteItemWithOthers(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"Important Item Don't Delete\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/delete_item/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	assertCount(t, db, "wishlist_items", 1)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=2").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "Important Item Don't Delete", "Wrong item_name")
}

func TestDeleteItemMultiple(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"secondaryItem\");")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"Important Item Don't Delete\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/delete_item/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	req = createFormRequest("/wishlist/delete_item/2", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	assertCount(t, db, "wishlist_items", 1)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=3").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "Important Item Don't Delete", "Wrong item_name")
}

func TestDeleteItemInvalidID(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/delete_item/5", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusNotFound, res.StatusCode)
	}

	assertCount(t, db, "wishlist_items", 1)

	var gotWishlistId int
	var gotName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&gotWishlistId, &gotName)
	if err != nil {
		t.Errorf("Error querying database for changed item: %v", err)
	}

	assert(t, gotWishlistId, 1, "Wrong wishlist_id")
	assert(t, gotName, "initialItem", "Wrong item_name")
}

func TestDeleteItemDoubleDelete(t *testing.T) {
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO wishlist_items (wishlist_id, item_name) VALUES (1, \"initialItem\");")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := createFormRequest("/wishlist/delete_item/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	assertRedirect(t, res, "/wishlist")

	req = createFormRequest("/wishlist/delete_item/1", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusNotFound, res.StatusCode)
	}

	assertCount(t, db, "wishlist_items", 0)
}
