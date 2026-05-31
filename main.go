package main

import (
	"flag"
	"log"

	"github.com/moltenwolfcub/giftsWebApp/app"
	"github.com/moltenwolfcub/giftsWebApp/config"
	"github.com/moltenwolfcub/giftsWebApp/database"
)

var migrateFlag bool

func main() {
	flag.BoolVar(&migrateFlag, "migrate", false, "Whether the database should apply any new migrations")
	flag.Parse()

	cfg := config.New()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	if migrateFlag {
		err := database.MigrateDB(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	app := app.NewAppServer(db, cfg)
	log.Fatal(app.Run())
}
