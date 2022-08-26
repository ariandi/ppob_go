// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"
)

type Querier interface {
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreatePartner(ctx context.Context, arg CreatePartnerParams) (Partner, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateProvider(ctx context.Context, arg CreateProviderParams) (Provider, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateRoleUser(ctx context.Context, arg CreateRoleUserParams) (RoleUser, error)
	CreateSelling(ctx context.Context, arg CreateSellingParams) (Selling, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteCategories(ctx context.Context, id int64) error
	DeletePartner(ctx context.Context, id int64) error
	DeleteProduct(ctx context.Context, id int64) error
	DeleteProvider(ctx context.Context, id int64) error
	DeleteRole(ctx context.Context, id int64) error
	DeleteRoleUser(ctx context.Context, id int64) error
	DeleteSelling(ctx context.Context, id int64) error
	DeleteTransaction(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetCategory(ctx context.Context, id int64) (Category, error)
	GetPartner(ctx context.Context, id int64) (Partner, error)
	GetPartnerByParams(ctx context.Context, arg GetPartnerByParamsParams) (Partner, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	GetProvider(ctx context.Context, id int64) (Provider, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	GetRoleUserByID(ctx context.Context, id int64) (RoleUser, error)
	GetRoleUserByRoleID(ctx context.Context, arg GetRoleUserByRoleIDParams) ([]RoleUser, error)
	GetRoleUserByUserID(ctx context.Context, arg GetRoleUserByUserIDParams) ([]RoleUser, error)
	GetSelling(ctx context.Context, id int64) (Selling, error)
	GetTransaction(ctx context.Context, id int64) (Transaction, error)
	GetTransactionByRefID(ctx context.Context, arg GetTransactionByRefIDParams) (Transaction, error)
	GetTransactionByTxID(ctx context.Context, txID string) (Transaction, error)
	GetTransactionPending(ctx context.Context, billID string) (Transaction, error)
	GetUser(ctx context.Context, id int64) (GetUserRow, error)
	GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error)
	ListCategory(ctx context.Context, arg ListCategoryParams) ([]Category, error)
	ListPartner(ctx context.Context, arg ListPartnerParams) ([]Partner, error)
	ListProduct(ctx context.Context, arg ListProductParams) ([]Product, error)
	ListProductByCatID(ctx context.Context, arg ListProductByCatIDParams) ([]Product, error)
	ListProvider(ctx context.Context, arg ListProviderParams) ([]Provider, error)
	ListRole(ctx context.Context, arg ListRoleParams) ([]Role, error)
	ListRoleUser(ctx context.Context, arg ListRoleUserParams) ([]RoleUser, error)
	ListRoleWithDelete(ctx context.Context, arg ListRoleWithDeleteParams) ([]Role, error)
	ListSelling(ctx context.Context, arg ListSellingParams) ([]Selling, error)
	ListSellingByParams(ctx context.Context, arg ListSellingByParamsParams) ([]Selling, error)
	ListTransaction(ctx context.Context, arg ListTransactionParams) ([]Transaction, error)
	ListUser(ctx context.Context, arg ListUserParams) ([]ListUserRow, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	UpdateInactiveCategory(ctx context.Context, arg UpdateInactiveCategoryParams) (Category, error)
	UpdateInactivePartner(ctx context.Context, arg UpdateInactivePartnerParams) (Partner, error)
	UpdateInactiveProduct(ctx context.Context, arg UpdateInactiveProductParams) (Product, error)
	UpdateInactiveProvider(ctx context.Context, arg UpdateInactiveProviderParams) (Provider, error)
	UpdateInactiveRole(ctx context.Context, arg UpdateInactiveRoleParams) (Role, error)
	UpdateInactiveRoleUser(ctx context.Context, arg UpdateInactiveRoleUserParams) (RoleUser, error)
	UpdateInactiveSelling(ctx context.Context, arg UpdateInactiveSellingParams) (Selling, error)
	UpdateInactiveTransaction(ctx context.Context, arg UpdateInactiveTransactionParams) (Transaction, error)
	UpdateInactiveUser(ctx context.Context, arg UpdateInactiveUserParams) (User, error)
	UpdatePartner(ctx context.Context, arg UpdatePartnerParams) (Partner, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateProvider(ctx context.Context, arg UpdateProviderParams) (Provider, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateRoleUser(ctx context.Context, arg UpdateRoleUserParams) (RoleUser, error)
	UpdateSelling(ctx context.Context, arg UpdateSellingParams) (Selling, error)
	UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) (Transaction, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
