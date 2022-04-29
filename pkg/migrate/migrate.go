package migrate

import (
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
)

const MIGRATIONSDIR = "./migrations"

func Migrate(db *sqlx.DB) {
	files, err := ioutil.ReadDir(MIGRATIONSDIR)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		file, err := ioutil.ReadFile(MIGRATIONSDIR + "/" + f.Name())

		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(string(file))

		if err != nil {
			log.Fatal(err)
		}
	}
}
