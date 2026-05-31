package main

import (
	"github.com/moltenwolfcub/giftsWebApp/app"
	"github.com/moltenwolfcub/giftsWebApp/config"
	"github.com/moltenwolfcub/giftsWebApp/database"

	_ "modernc.org/sqlite"
)

const MIGRATE = true

func main() {
	cfg := config.New()

	db := database.ConnectDB()
	if MIGRATE {
		database.MigrateDB(db)
	}

	app := app.NewAppServer(db, cfg)
	app.Run()
}
