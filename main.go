package main

import (
	"github.com/moltenwolfcub/giftsWebApp/app"
	"github.com/moltenwolfcub/giftsWebApp/config"
	"github.com/moltenwolfcub/giftsWebApp/database"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.New()

	db := database.LoadDatabase()

	app := app.NewAppServer(db, cfg)
	app.Run()
}
