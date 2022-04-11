package migrate

import (
	"database/sql"
	"io/ioutil"
	"log"
)

const MIGRATIONSDIR = "./migrations"

func Migrate(db *sql.DB) {
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
