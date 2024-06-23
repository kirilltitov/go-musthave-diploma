package storage

import (
	"context"
	"embed"
	"io/fs"

	"github.com/jackc/tern/v2/migrate"
	"github.com/pkg/errors"

	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (s PgSQL) InitDB(ctx context.Context) error {
	conn, err := s.Conn.Acquire(ctx)
	if err != nil {
		return err
	}

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), "schema_version")
	if err != nil {
		return errors.Wrap(err, "Unable to create a migrator")
	}

	fsys, err := fs.Sub(embedMigrations, "migrations") // без этого не работает, и я не пойму, почему
	if err != nil {
		return errors.Wrap(err, "Unable load embed migrations")
	}

	err = migrator.LoadMigrations(fsys)
	if err != nil {
		return errors.Wrap(err, "Unable to load migrations")
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to migrate")
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to get current schema version")
	}

	utils.Log.Infof("Migration done. Current schema version: %v\n", ver)

	return nil
}
