package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/shanmukhsista/go-graphql-starter/pkg/common/config"
)

const (
	kErrUnableToFetchTransaction = "kErrUnableToFetchTransaction"
)

type Database struct {
	Conn *pgxpool.Pool
}

// TransactionManager runs logic inside a single database transaction
type TransactionManager interface {
	// WithinTransaction runs a function within a database transaction.
	//
	// Transaction is propagated in the context,
	// so it is important to propagate it to underlying repositories.
	// Function commits if error is nil, and rollbacks if not.
	// It returns the same error.
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

func ProvideNewDatabaseConnection(pgpool *pgxpool.Pool) *Database {
	return &Database{Conn: pgpool}
}

func ProvideNewPostgresTransactor(db *Database) (TransactionManager, error) {
	return db, nil
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context and creates a new tx if it doesn't exist.
func (db *Database) ExtractTx(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txKey{}).(pgx.Tx)
	return tx
}

func (db *Database) GetExistingOrNewTransaction(ctx context.Context) (pgx.Tx, error) {
	tx := db.ExtractTx(ctx)
	if tx == nil {
		tx, err := db.Conn.Begin(ctx)
		if err != nil {
			return nil, err
		}
		return tx, nil
	}
	return tx, nil
}

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error
func (db *Database) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	// begin transaction
	tx, err := db.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			log.Error().Err(errRollback)
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit(ctx); errCommit != nil {
		log.Error().Err(errCommit)
		return errCommit
	}
	return nil
}

func (db *Database) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	tx := db.ExtractTxWithoutError(ctx)
	var fetchedRow pgx.Row
	if tx != nil {
		// execute in a transaction
		fetchedRow = tx.QueryRow(ctx, sql, args...)
	} else {
		// execute independently.
		fetchedRow = db.Conn.QueryRow(ctx, sql, args...)
	}
	return fetchedRow
}

func (db *Database) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	tx := db.ExtractTxWithoutError(ctx)
	var rows pgx.Rows
	var err error
	if tx != nil {
		// execute in a transaction
		rows, err = tx.Query(ctx, sql, args...)
	} else {
		// execute independently.
		rows, err = db.Conn.Query(ctx, sql, args...)
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *Database) ExtractTxWithoutError(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txKey{}).(pgx.Tx)

	return tx
}

func (db *Database) ExtractTxWithError(ctx context.Context) (pgx.Tx, error) {
	tx, _ := ctx.Value(txKey{}).(pgx.Tx)
	if tx == nil {
		return nil, errors.New(kErrUnableToFetchTransaction)
	}
	return tx, nil
}

func ProvidePgConnectionPool() (*pgxpool.Pool, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbUrl := config.MustGetString("db.postgres.url")
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ProvideNewDatabase(pool *pgxpool.Pool) (*Database, error) {
	database := ProvideNewDatabaseConnection(pool)
	return database, nil
}
