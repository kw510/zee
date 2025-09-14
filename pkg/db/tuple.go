package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Tuple struct {
	Parent Set
	Child  Set
}

func (t Tuple) CreateTx(ctx context.Context, tx pgx.Tx) error {
	q := `
		INSERT INTO tuples (parent_namespace, parent_id, parent_relation, child_namespace, child_id, child_relation)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`
	_, err := tx.Exec(ctx, q, t.Parent.Namespace, t.Parent.Id, t.Parent.Relation, t.Child.Namespace, t.Child.Id, t.Child.Relation)
	if err != nil {
		return fmt.Errorf("failed to insert tuple: %w", err)
	}
	return nil
}

func (t Tuple) DeleteTx(ctx context.Context, tx pgx.Tx) error {
	q := `
		DELETE FROM tuples
		WHERE parent_namespace = $1 AND parent_id = $2 AND parent_relation = $3 AND child_namespace = $4 AND child_id = $5 AND child_relation = $6
	`
	_, err := tx.Exec(ctx, q, t.Parent.Namespace, t.Parent.Id, t.Parent.Relation, t.Child.Namespace, t.Child.Id, t.Child.Relation)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}
	return nil
}

func (t Tuple) Exists(ctx context.Context) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT 1 FROM tuples
			WHERE parent_namespace = $1 AND parent_id = $2 AND parent_relation = $3 AND child_namespace = $4 AND child_id = $5 AND child_relation = $6
		)
	`
	var exists bool
	err := pg.QueryRow(ctx, q, t.Parent.Namespace, t.Parent.Id, t.Parent.Relation, t.Child.Namespace, t.Child.Id, t.Child.Relation).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if tuple exists: %w", err)
	}
	return exists, nil
}
