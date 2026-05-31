package controller_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moltenwolfcub/giftsWebApp/controller"
)

func genDataBase() (*sql.DB, func()) {
	return nil, nil
}

func TestAddItem(t *testing.T) {
	db, teardownDB := genDataBase()
	defer teardownDB()
	// cfg := config.New()
	router := controller.BuildRouter(db)

	req := httptest.NewRequest("POST", "/wishlist/add_item", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code: %d", res.StatusCode)
	}

	// db.QueryRow()
	// check row correctly added to wishlist_items
}
