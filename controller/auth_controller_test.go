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
			"foo);DROP TABLE users;",
			"foo);DROP TABLE users;",
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
				t.Errorf("Error querying database for registered user: %v", err)
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

func TestRegisterMultiple(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "John Smith",
		"password":         "password",
		"confirm-password": "password",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/")

	req = testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "Eva Smith",
		"password":         "password",
		"confirm-password": "password",
	})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/")

	testhelpers.AssertCount(t, db, "users", 2)

	userExists, err := models.UsernameExists(db, "John Smith")
	if err != nil {
		t.Errorf("Error testing for username existance: %v", err)
	}
	if !userExists {
		t.Errorf("Registered user doesn't exist in database")
	}
	userExists, err = models.UsernameExists(db, "Eva Smith")
	if err != nil {
		t.Errorf("Error testing for username existance: %v", err)
	}
	if !userExists {
		t.Errorf("Registered user doesn't exist in database")
	}

	var gotUsername1, gotPassword1, gotSalt1 string
	err = db.QueryRow("SELECT username, password, salt FROM users WHERE id=1").Scan(&gotUsername1, &gotPassword1, &gotSalt1)
	if err != nil {
		t.Errorf("Error querying database for registered user: %v", err)
	}

	testhelpers.Assert(t, gotUsername1, "John Smith", "Incorrect username stored in database")

	normalPassword, err := hex.DecodeString(gotPassword1)
	if err != nil {
		t.Errorf("Error decoding password from hex: %v", err)
	}
	normalSalt, err := hex.DecodeString(gotSalt1)
	if err != nil {
		t.Errorf("Error decoding salt from hex: %v", err)
	}

	expectedPassword := models.HashPasswordWithSalt("password", cfg, normalSalt)
	testhelpers.Assert(t, string(normalPassword), string(expectedPassword), "wrong password hash")

	var gotUsername2, gotPassword2, gotSalt2 string
	err = db.QueryRow("SELECT username, password, salt FROM users WHERE id=2").Scan(&gotUsername2, &gotPassword2, &gotSalt2)
	if err != nil {
		t.Errorf("Error querying database for registered user: %v", err)
	}

	testhelpers.Assert(t, gotUsername2, "Eva Smith", "Incorrect username stored in database")

	normalPassword, err = hex.DecodeString(gotPassword2)
	if err != nil {
		t.Errorf("Error decoding password from hex: %v", err)
	}
	normalSalt, err = hex.DecodeString(gotSalt2)
	if err != nil {
		t.Errorf("Error decoding salt from hex: %v", err)
	}

	expectedPassword = models.HashPasswordWithSalt("password", cfg, normalSalt)
	testhelpers.Assert(t, string(normalPassword), string(expectedPassword), "wrong password hash")
}

func TestRegisterEmptyUsername(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "",
		"password":         "password",
		"confirm-password": "password",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/auth/register")

	testhelpers.AssertCount(t, db, "users", 0)
}

func TestRegisterDuplicateUsername(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")
	db.Exec("INSERT INTO users (username, password, salt) VALUES (\"takenName\",\"0\",\"0\");")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "takenName",
		"password":         "password",
		"confirm-password": "password",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/auth/register")

	testhelpers.AssertCount(t, db, "users", 1)
}

func TestRegisterEmptyPassword(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "name",
		"password":         "",
		"confirm-password": "",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/auth/register")

	testhelpers.AssertCount(t, db, "users", 0)
}

func TestRegisterBadConfirmPassword(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "name",
		"password":         "something",
		"confirm-password": "a different thing",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/auth/register")

	testhelpers.AssertCount(t, db, "users", 0)
}

func TestRegisterPasswordUniqueHash(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "John Smith",
		"password":         "password123",
		"confirm-password": "password123",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/")

	req = testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "Eva Smith",
		"password":         "password123",
		"confirm-password": "password123",
	})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res = w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/")

	testhelpers.AssertCount(t, db, "users", 2)

	var johnHash string
	err := db.QueryRow("SELECT password FROM users WHERE id=1").Scan(&johnHash)
	if err != nil {
		t.Errorf("Error querying database for registerd user: %v", err)
	}
	var evaHash string
	err = db.QueryRow("SELECT password FROM users WHERE id=2").Scan(&evaHash)
	if err != nil {
		t.Errorf("Error querying database for registerd user: %v", err)
	}

	if evaHash == johnHash {
		t.Error("The same password from 2 useres has the same hash")
	}
}

func TestRegisterNoBody(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	req := httptest.NewRequest("POST", "/auth/register", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/auth/register")

	testhelpers.AssertCount(t, db, "users", 0)
}

func TestLogin(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	//create user
	req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "user1",
		"password":         "password",
		"confirm-password": "password",
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//test user
	req = testhelpers.CreateFormRequest("/auth/login", map[string]string{
		"username": "user1",
		"password": "password",
	})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/")
	// when login is more complex check cookie that user is in fact authenticated
}

func TestLoginWrongPassword(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	//create user
	req := testhelpers.CreateFormRequest("/auth/register", map[string]string{
		"username":         "user1",
		"password":         "password",
		"confirm-password": "password",
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//test user
	req = testhelpers.CreateFormRequest("/auth/login", map[string]string{
		"username": "user1",
		"password": "wrongPassword",
	})
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/auth/login")
}

func TestLoginInvalidUser(t *testing.T) {
	db := testhelpers.GenTestDatabase(t)
	db.Exec("INSERT INTO wishlists DEFAULT VALUES")

	cfg := config.New()
	router := controller.BuildRouter(db, cfg)

	//test user
	req := testhelpers.CreateFormRequest("/auth/login", map[string]string{
		"username": "doesntExist",
		"password": "pass",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	t.Cleanup(func() {
		res.Body.Close()
	})

	testhelpers.AssertRedirect(t, res, "/auth/login")
}
