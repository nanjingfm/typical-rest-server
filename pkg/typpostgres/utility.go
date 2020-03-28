package typpostgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang-migrate/migrate"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

// PostgresUtility is postgres utility
type PostgresUtility struct{}

// Utility return new instance of PostgresUtility
func Utility() *PostgresUtility {
	return &PostgresUtility{}
}

// Commands of module
func (*PostgresUtility) Commands(c *typbuildtool.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres Database Tool",
			Subcommands: []*cli.Command{
				{
					Name:  "create",
					Usage: "Create New Database",
					Action: func(cliCtx *cli.Context) (err error) {
						return create(c)
					},
				},
				{
					Name:  "drop",
					Usage: "Drop Database",
					Action: func(cliCtx *cli.Context) (err error) {
						return drop(c)
					},
				},
				{
					Name:  "migrate",
					Usage: "Migrate Database",
					Action: func(cliCtx *cli.Context) (err error) {
						return migrateDB(c)
					},
				},
				{
					Name:  "rollback",
					Usage: "Rollback Database",
					Action: func(cliCtx *cli.Context) (err error) {
						return rollbackDB(c)
					},
				},
				{
					Name:  "seed",
					Usage: "Data seeding",
					Action: func(cliCtx *cli.Context) (err error) {
						return seed(c)
					},
				},
				{
					Name:  "reset",
					Usage: "Reset Database",
					Action: func(cliCtx *cli.Context) (err error) {
						return reset(c)
					},
				},
				{
					Name:  "console",
					Usage: "PostgreSQL Interactive",
					Action: func(cliCtx *cli.Context) (err error) {
						return console(c)
					},
				},
			},
		},
	}
}

func reset(c *typbuildtool.Context) (err error) {
	if err = drop(c); err != nil {
		return
	}
	if err = create(c); err != nil {
		return
	}
	if err = migrateDB(c); err != nil {
		return
	}
	if err = seed(c); err != nil {
		return
	}
	return
}

func create(c *typbuildtool.Context) (err error) {
	var conn *sql.DB
	var cfg *Config

	if cfg, err = RetrieveConfig(c); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		return
	}

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	c.Infof("Postgres: %s", query)
	_, err = conn.Exec(query)
	return
}

func drop(c *typbuildtool.Context) (err error) {
	var conn *sql.DB
	var cfg *Config

	if cfg, err = RetrieveConfig(c); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	c.Infof("Postgres: %s", query)
	_, err = conn.Exec(query)
	return
}

func console(c *typbuildtool.Context) (err error) {
	var cfg *Config
	if cfg, err = RetrieveConfig(c); err != nil {
		return
	}

	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func migrateDB(c *typbuildtool.Context) (err error) {
	var (
		migration *migrate.Migrate
		sourceURL = "file://" + defaultMigrationSource
		cfg       *Config
	)

	if cfg, err = RetrieveConfig(c); err != nil {
		return
	}

	c.Infof("Migrate database from source '%s'", sourceURL)
	if migration, err = migrate.New(sourceURL, dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func rollbackDB(c *typbuildtool.Context) (err error) {
	var migration *migrate.Migrate
	var cfg *Config
	if cfg, err = RetrieveConfig(c); err != nil {
		return
	}

	sourceURL := "file://" + defaultMigrationSource
	c.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func seed(c *typbuildtool.Context) (err error) {
	var conn *sql.DB
	var cfg *Config
	if cfg, err = RetrieveConfig(c); err != nil {
		return
	}
	if conn, err = sql.Open("postgres", dataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	files, _ := ioutil.ReadDir(defaultSeedSource)
	for _, f := range files {
		sqlFile := defaultSeedSource + "/" + f.Name()
		c.Infof("Execute seed '%s'", sqlFile)
		var b []byte
		if b, err = ioutil.ReadFile(sqlFile); err != nil {
			c.Warn(err.Error())
			continue
		}
		if _, err = conn.Exec(string(b)); err != nil {
			c.Warn(err.Error())
		}
	}
	return
}