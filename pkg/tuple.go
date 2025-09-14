package pkg

import (
	"context"

	"github.com/kw510/z/pkg/db"
)

func Check(ctx context.Context, t db.Tuple) (bool, error) {
	// Check if the tuple exists directly
	exists, err := t.Exists(ctx)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	// Check if the tuple exists in any of the parent's subsets
	ss, err := t.Parent.Subsets(ctx)
	if err != nil {
		return false, err
	}

	for _, s := range ss {
		exists, err := Check(ctx, db.Tuple{Parent: s, Child: t.Child})
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}
	}
	return false, nil
}
