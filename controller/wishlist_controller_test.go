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
	db := genDataBase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	// cfg := config.New()
	router := controller.BuildRouter(db)

	testName := "Test Item"

	formContents := url.Values{}
	formContents.Set("itemName", testName)

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

	var wishlistId int
	var itemName string
	err := db.QueryRow("SELECT wishlist_id, item_name FROM wishlist_items WHERE id=1").Scan(&wishlistId, &itemName)
	if err != nil {
		t.Errorf("Error querying database for inserted item: %v", err)
	}

	if wishlistId != 1 {
		t.Errorf("Wrong wishlist_id. Expected: %d, Got %d", 1, wishlistId)
	}
	if itemName != testName {
		t.Errorf("Wrong item_name. Expected: %s, Got %s", testName, itemName)
	}
}
