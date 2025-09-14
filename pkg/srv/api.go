package srv

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"github.com/kw510/z/pkg/db"
	apiv1 "github.com/kw510/z/pkg/gen/z/api/v1"
)

type ApiServer struct{}

func (srv *ApiServer) Check(ctx context.Context, req *connect.Request[apiv1.CheckRequest]) (*connect.Response[apiv1.CheckResponse], error) {
	t := db.Tuple{
		Parent: db.Set{
			Namespace: req.Msg.GetParent().GetNamespace(),
			Id:        req.Msg.GetParent().GetId(),
			Relation:  req.Msg.GetParent().GetRelation(),
		},
		Child: db.Set{
			Namespace: req.Msg.GetChild().GetNamespace(),
			Id:        req.Msg.GetChild().GetId(),
			Relation:  req.Msg.GetChild().GetRelation(),
		},
	}

	// Check if the tuple exists directly
	exists, err := t.Exists(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if exists {
		return &connect.Response[apiv1.CheckResponse]{
			Msg: &apiv1.CheckResponse{
				Ok: true,
			},
		}, nil
	}

	// Check if the tuple exists in any of the parent's subsets
	ss, err := t.Parent.Subsets(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, s := range ss {
		res, err := srv.Check(ctx, &connect.Request[apiv1.CheckRequest]{
			Msg: &apiv1.CheckRequest{
				Parent: &apiv1.Set{
					Namespace: s.Namespace,
					Id:        s.Id,
					Relation:  s.Relation,
				},
				Child: &apiv1.Set{
					Namespace: t.Child.Namespace,
					Id:        t.Child.Id,
					Relation:  t.Child.Relation,
				},
			},
		})
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if res.Msg.Ok {
			return &connect.Response[apiv1.CheckResponse]{
				Msg: &apiv1.CheckResponse{
					Ok: true,
				},
			}, nil
		}
	}
	return &connect.Response[apiv1.CheckResponse]{
		Msg: &apiv1.CheckResponse{
			Ok: false,
		},
	}, nil
}

func (srv *ApiServer) Write(ctx context.Context, req *connect.Request[apiv1.WriteRequest]) (*connect.Response[apiv1.WriteResponse], error) {
	err := db.Tx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		for _, tuple := range req.Msg.GetAdds() {
			t := db.Tuple{
				Parent: db.Set{
					Namespace: tuple.GetParent().GetNamespace(),
					Id:        tuple.GetParent().GetId(),
					Relation:  tuple.GetParent().GetRelation(),
				},
				Child: db.Set{
					Namespace: tuple.GetChild().GetNamespace(),
					Id:        tuple.GetChild().GetId(),
					Relation:  tuple.GetChild().GetRelation(),
				},
			}
			err := t.CreateTx(ctx, tx)
			if err != nil {
				return fmt.Errorf("failed to create tuple: %w", err)
			}
		}
		for _, tuple := range req.Msg.GetRemoves() {
			t := db.Tuple{
				Parent: db.Set{
					Namespace: tuple.GetParent().GetNamespace(),
					Id:        tuple.GetParent().GetId(),
					Relation:  tuple.GetParent().GetRelation(),
				},
				Child: db.Set{
					Namespace: tuple.GetChild().GetNamespace(),
					Id:        tuple.GetChild().GetId(),
					Relation:  tuple.GetChild().GetRelation(),
				},
			}
			err := t.DeleteTx(ctx, tx)
			if err != nil {
				return fmt.Errorf("failed to delete tuple: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &connect.Response[apiv1.WriteResponse]{}, nil
}
