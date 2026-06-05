package testhelpers

import (
	"database/sql"
	"testing"

	"github.com/moltenwolfcub/giftsWebApp/database"
)

func GenTestDatabase(t *testing.T) *sql.DB {
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
