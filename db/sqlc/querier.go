// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"
)

type Querier interface {
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteRole(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetRole(ctx context.Context, id int64) (Role, error)
	GetUser(ctx context.Context, id int64) (User, error)
	ListRole(ctx context.Context, arg ListRoleParams) ([]Role, error)
	ListUser(ctx context.Context, arg ListUserParams) ([]User, error)
	UpdateInactiveRole(ctx context.Context, arg UpdateInactiveRoleParams) (Role, error)
	UpdateInactiveUser(ctx context.Context, arg UpdateInactiveUserParams) (User, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
