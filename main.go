package main

import (
	"log"

	"github.com/moltenwolfcub/giftsWebApp/app"
	"github.com/moltenwolfcub/giftsWebApp/config"
	"github.com/moltenwolfcub/giftsWebApp/database"

	_ "modernc.org/sqlite"
)

const MIGRATE = true

func main() {
	cfg := config.New()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	if MIGRATE {
		err := database.MigrateDB(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	app := app.NewAppServer(db, cfg)
	log.Fatal(app.Run())
}
