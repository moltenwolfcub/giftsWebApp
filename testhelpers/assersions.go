package testhelpers

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"
)

func AssertRedirect(t *testing.T, res *http.Response, expected string) {
	if res.StatusCode != http.StatusFound { //TODO check all status codes used are correct
		t.Errorf("Wrong http status code. Expected: %d, Got: %d", http.StatusFound, res.StatusCode)
	}

	if res.Header["Location"][0] != expected {
		t.Errorf("Redirected to wrong webpage. Expected: %s, Got %s", expected, res.Header["Location"][0])
	}
}

func AssertCount(t *testing.T, db *sql.DB, table string, expected int) {
	var found int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&found)
	if err != nil {
		t.Errorf("Error counting rows in database: %v", err)
	}

	Assert(t, found, expected, "Wrong number of items in "+table)
}

func Assert[K comparable](t *testing.T, got, want K, msg string) {
	if got != want {
		t.Errorf("%s. Expected: %v, Got %v", msg, want, got)
	}

}
