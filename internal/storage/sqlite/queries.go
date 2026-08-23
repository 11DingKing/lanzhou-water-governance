package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func WithTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err = fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%w; rollback: %v", err, rollbackErr)
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
func IsConstraint(err error) bool {
	return err != nil && (errors.Is(err, sql.ErrNoRows) || containsConstraint(err.Error()))
}
func containsConstraint(message string) bool {
	return len(message) > 0 && (message[0] == 'U' || message[0] == 'u' || message[0] == 'C' || message[0] == 'c')
}
