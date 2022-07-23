package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/sirupsen/logrus"
)

type Store interface {
	Querier
	//TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, req dto.CreateUserRequest, authPayload *token.Payload) (dto.UserResponse, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *SQLStore) CreateUserTx(ctx context.Context, req dto.CreateUserRequest, authPayload *token.Payload) (dto.UserResponse, error) {
	logrus.Println("[Store CreateUserTx] start.")
	var result dto.UserResponse
	err := store.execTransaction(ctx, func(q *Queries) error {
		var err error
		userPayload, err := q.GetUserByUsername(ctx, authPayload.Username)
		user, err := q.CreateUser(ctx, CreateUserParams{
			Name:           req.Name,
			Email:          req.Email,
			Username:       req.Username,
			Password:       sql.NullString{String: req.Password, Valid: true},
			Balance:        sql.NullString{String: "0.00", Valid: true},
			Phone:          req.Phone,
			IdentityNumber: req.IdentityNumber,
			CreatedBy:      sql.NullInt64{Int64: userPayload.ID, Valid: true},
		})
		if err != nil {
			return err
		}

		arg := CreateRoleUserParams{
			RoleID:    req.RoleID,
			UserID:    user.ID,
			CreatedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
		}
		roleUser, err := q.CreateRoleUser(ctx, arg)
		if err != nil {
			return err
		}

		result.Name = user.Name
		result.ID = user.ID
		result.Phone = user.Phone
		result.Username = user.Username
		result.Email = user.Email
		result.Balance = user.Balance
		result.IdentityNumber = user.IdentityNumber

		resRoleUser := dto.RoleUser{
			ID:        roleUser.ID,
			RoleID:    roleUser.RoleID,
			UserID:    roleUser.UserID,
			CreatedAt: roleUser.CreatedAt,
			UpdatedAt: roleUser.UpdatedAt,
			CreatedBy: roleUser.CreatedBy,
			UpdatedBy: roleUser.UpdatedBy,
		}

		result.Role = append(result.Role, resRoleUser)

		return nil
	})

	return result, err
}
