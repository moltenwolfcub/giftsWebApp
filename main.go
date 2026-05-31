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

	db := database.LoadDatabase(MIGRATE)

	app := app.NewAppServer(db, cfg)
	app.Run()
}
