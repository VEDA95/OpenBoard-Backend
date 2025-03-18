package db

import (
	"VEDA95/open_board/api/internal/log"
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	zerologadapter "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"os"
)

type DB struct {
	db      *pgxpool.Pool
	context context.Context
}

var Instance *DB

func CreateDBInstance(dbUrl string, logger zerolog.Logger) (*DB, error) {
	connConfig, err := pgxpool.ParseConfig(dbUrl)

	if err != nil {
		return nil, err
	}

	connConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   zerologadapter.NewLogger(logger),
		LogLevel: tracelog.LogLevelDebug,
	}
	connConfig.AfterConnect = func(ctx context.Context, connection *pgx.Conn) error {
		pgxuuid.Register(connection.TypeMap())
		return nil
	}
	dbContext := context.Background()
	db, err := pgxpool.NewWithConfig(dbContext, connConfig)

	if err != nil {
		return nil, err
	}

	return &DB{db: db, context: dbContext}, nil
}

func InitializeDBInstance() error {
	dbUrl := os.Getenv("DATABASE_URL")

	if len(dbUrl) == 0 {
		return errors.New("DB_URL has not been set")
	}

	instance, err := CreateDBInstance(dbUrl, log.Logger)

	if err != nil {
		return err
	}

	Instance = instance

	return nil
}

func (store *DB) One(builder sqlbuilder.Builder, output interface{}) error {
	query, args := builder.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := store.db.Query(store.context, query, args...)

	if err != nil {
		return err
	}

	if err := pgxscan.ScanOne(output, rows); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		return err
	}

	return nil
}

func (store *DB) Many(builder sqlbuilder.Builder, output interface{}) error {
	query, args := builder.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := store.db.Query(store.context, query, args...)

	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(output, rows); err != nil {
		return err
	}

	return nil
}

func (store *DB) Exec(builder sqlbuilder.Builder) error {
	query, args := builder.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := store.db.Exec(store.context, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (store *DB) Close() {
	store.db.Close()
}

func (store *DB) Begin() (*Transaction, error) {
	tx, err := store.db.Begin(store.context)

	if err != nil {
		return nil, err
	}

	return &Transaction{tx: tx, context: context.Background()}, nil
}
