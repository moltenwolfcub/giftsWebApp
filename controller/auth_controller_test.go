package controller_test

import (
	"encoding/hex"
	"net/http/httptest"
	"testing"

	"github.com/moltenwolfcub/giftsWebApp/config"
	"github.com/moltenwolfcub/giftsWebApp/controller"
	"github.com/moltenwolfcub/giftsWebApp/models"
	"github.com/moltenwolfcub/giftsWebApp/testhelpers"
)

func TestRegister(t *testing.T) {
	var tests = []struct {
		testName string
		username string
		password string
	}{
		{
			"Normal",
			"John Smith",
			"password",
		},
		{
			"Numeric Name",
			"12345",
			"password",
		},
		{
			"Unicode",
			"!#$%& ̄f ̅ ̆ ̇ ̈ ̉ ̊ ̋ ̌ ̍ ̎ ̏ ̐ ̑ ̒ ̓ ̔ ̕ ζδψφΔΓΨΣΩ",
			"!#$%& ̄f ̅ ̆ ̇ ̈ ̉ ̊ ̋ ̌ ̍ ̎ ̏ ̐ ̑ ̒ ̓ ̔ ̕ ζδψφΔΓΨΣΩ",
		},
		{
			"SQL Injection",
			"foo);DROP TABLE wishlist_items;",
			"foo);DROP TABLE wishlist_items;",
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.testName, func(t *testing.T) {
			db := testhelpers.GenTestDatabase(t)
			db.Exec("INSERT INTO wishlists DEFAULT VALUES")

			cfg := config.New()
			router := controller.BuildRouter(db, cfg)

			req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
				"username":         testcase.username,
				"password":         testcase.password,
				"confirm-password": testcase.password,
			})
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			res := w.Result()
			t.Cleanup(func() {
				res.Body.Close()
			})

			testhelpers.AssertRedirect(t, res, "/")

			testhelpers.AssertCount(t, db, "users", 1)

			userExists, err := models.UsernameExists(db, testcase.username)
			if err != nil {
				t.Errorf("Error testing for username existance: %v", err)
			}
			if !userExists {
				t.Errorf("Registered user doesn't exist in database")
			}

			var gotUsername, gotPassword, gotSalt string
			err = db.QueryRow("SELECT username, password, salt FROM users WHERE id=1").Scan(&gotUsername, &gotPassword, &gotSalt)
			if err != nil {
				t.Errorf("Error querying database for inserted item: %v", err)
			}

			testhelpers.Assert(t, gotUsername, testcase.username, "Incorrect username stored in database")

			normalPassword, err := hex.DecodeString(gotPassword)
			if err != nil {
				t.Errorf("Error decoding password from hex: %v", err)
			}
			normalSalt, err := hex.DecodeString(gotSalt)
			if err != nil {
				t.Errorf("Error decoding salt from hex: %v", err)
			}

			expectedPassword := models.HashPasswordWithSalt(testcase.password, cfg, normalSalt)
			testhelpers.Assert(t, string(normalPassword), string(expectedPassword), "wrong password hash")
		})
	}
}
