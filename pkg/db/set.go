package db

import (
	"context"
	"fmt"
)

type Set struct {
	Namespace string
	Id        string
	Relation  string
}

func (s Set) Children(ctx context.Context) ([]Set, error) {
	q := `
		SELECT child_namespace, child_id, child_relation FROM tuples WHERE parent_namespace = $1 AND parent_id = $2 AND parent_relation = $3
	`
	rows, err := pg.Query(ctx, q, s.Namespace, s.Id, s.Relation)
	if err != nil {
		return nil, fmt.Errorf("failed to list children: %w", err)
	}
	defer rows.Close()

	children := []Set{}
	for rows.Next() {
		var child Set
		err := rows.Scan(&child.Namespace, &child.Id, &child.Relation)
		if err != nil {
			return nil, fmt.Errorf("failed to scan child: %w", err)
		}
		children = append(children, child)
	}

	return children, nil
}

func (s Set) Subsets(ctx context.Context) ([]Set, error) {
	children, err := s.Children(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list children: %w", err)
	}

	out := []Set{}
	for _, child := range children {
		if child.Relation == "" {
			continue
		}
		out = append(out, child)
	}
	return out, nil
}
