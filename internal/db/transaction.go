package db

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	tx      pgx.Tx
	context context.Context
}

func (transaction *Transaction) One(builder sqlbuilder.Builder, output interface{}) error {
	query, args := builder.Build()
	rows, err := transaction.tx.Query(transaction.context, query, args...)

	if err != nil {
		return err
	}

	if err := pgxscan.ScanOne(output, rows); err != nil {
		return err
	}

	return nil
}

func (transaction *Transaction) Many(builder sqlbuilder.Builder, output interface{}) error {
	query, args := builder.Build()
	rows, err := transaction.tx.Query(transaction.context, query, args...)

	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(output, rows); err != nil {
		return err
	}

	return nil
}

func (transaction *Transaction) Exec(builder sqlbuilder.Builder) error {
	query, args := builder.Build()
	_, err := transaction.tx.Exec(transaction.context, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (transaction *Transaction) Commit() error {
	if err := transaction.tx.Commit(transaction.context); err != nil {
		if err := transaction.tx.Rollback(transaction.context); err != nil {
			return err
		}

		return err
	}

	return nil
}
