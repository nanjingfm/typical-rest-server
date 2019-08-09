package postgres

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/dbtool"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// Module for postgres
func Module() *typictx.Module {
	tool := dbtool.NewPostgresTool()

	return &typictx.Module{
		Name:      "postgres",
		ShortName: "pg",
		Usage:     "Postgres Database Module",
		SideEffects: []*typictx.SideEffect{
			typictx.NewSideEffect("github.com/lib/pq"),
			typictx.NewSideEffect("github.com/golang-migrate/migrate/database/postgres").ExcludeApp(),
			typictx.NewSideEffect("github.com/golang-migrate/migrate/source/file").ExcludeApp(),
		},
		Commands: []typictx.Command{
			{Name: "create", Usage: "Create New Database", ActionFunc: tool.CreateDB},
			{Name: "drop", Usage: "Drop Database", ActionFunc: tool.DropDB},
			{Name: "migrate", Usage: "Migrate Database", ActionFunc: tool.MigrateDB},
			{Name: "rollback", Usage: "Rollback Database", ActionFunc: tool.RollbackDB},
		},
		OpenFunc: func(cfg *Config) (*sql.DB, error) {
			log.Info("Open postgres connection")
			return sql.Open(cfg.DriverName(), cfg.DataSource())
		},
		CloseFunc: func(db *sql.DB) error {
			log.Info("Close postgres connection")
			return db.Close()
		},
	}
}