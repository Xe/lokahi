package database

import (
	"database/sql"

	"github.com/Xe/lokahi/internal/database/dmigrations"
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"
)

func Migrate(durl string) error {
	s := bindata.Resource(dmigrations.AssetNames(),
		func(name string) ([]byte, error) {
			return dmigrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", d, durl)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		return err
	}

	return nil
}

func Destroy(durl string) error {
	s := bindata.Resource(dmigrations.AssetNames(),
		func(name string) ([]byte, error) {
			return dmigrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", d, durl)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", durl)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Exec("DROP TABLE schema_migrations")

	return nil
}
