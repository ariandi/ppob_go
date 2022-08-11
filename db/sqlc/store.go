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
	CreateUserTx(ctx context.Context, req CreateUserParams, authPayload *token.Payload, roleId int64) (dto.UserResponse, error)
	UpdateUserTx(ctx context.Context, req UpdateUserParams, authPayload *token.Payload, roleId int64) (dto.UserResponse, error)
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

func (store *SQLStore) CreateUserTx(ctx context.Context, req CreateUserParams, authPayload *token.Payload, roleId int64) (dto.UserResponse, error) {
	logrus.Println("[Store CreateUserTx] start.")
	logrus.Println("[Store CreateUserTx] request is ", req)
	var result dto.UserResponse
	err := store.execTransaction(ctx, func(q *Queries) error {
		var err error
		userPayload, err := q.GetUserByUsername(ctx, authPayload.Username)
		user, err := q.CreateUser(ctx, CreateUserParams{
			Name:           req.Name,
			Email:          req.Email,
			Username:       req.Username,
			Password:       sql.NullString{String: req.Password.String, Valid: true},
			Balance:        sql.NullString{String: "0.00", Valid: true},
			Phone:          req.Phone,
			IdentityNumber: req.IdentityNumber,
			CreatedBy:      sql.NullInt64{Int64: userPayload.ID, Valid: true},
		})
		if err != nil {
			logrus.Println("[Store CreateUserTx] user is ", req.Email)
			logrus.Println("[Store CreateUserTx] error create user is ", err)
			return err
		}

		logrus.Println("[Store CreateUserTx] user ID is ", user.ID)
		arg := CreateRoleUserParams{
			RoleID:    roleId,
			UserID:    user.ID,
			CreatedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
		}
		roleUser, err := q.CreateRoleUser(ctx, arg)
		if err != nil {
			logrus.Println("[Store CreateUserTx] error create role is ", err)
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

func (store *SQLStore) UpdateUserTx(ctx context.Context, req UpdateUserParams, authPayload *token.Payload, roleId int64) (dto.UserResponse, error) {
	logrus.Println("[Store UpdateUserTx] start.")
	logrus.Println("[Store UpdateUserTx] request is ", req)
	var result dto.UserResponse
	err := store.execTransaction(ctx, func(q *Queries) error {
		var err error

		if req.Name != "" {
			req.SetName = true
		}
		if req.Phone != "" {
			req.SetPhone = true
		}
		if req.Password.Valid {
			req.SetPassword = true
		}
		if req.Email != "" {
			req.SetEmail = true
		}
		if req.BankCode.Valid {
			req.SetBankCode = true
		}
		if req.Balance.Valid {
			if authPayload.Username == "dbduabelas" {
				req.SetBalance = true
			} else {
				req.SetBalance = false
			}
		}

		userPayload, err := q.GetUserByUsername(ctx, authPayload.Username)
		req.UpdatedBy = sql.NullInt64{Int64: userPayload.ID, Valid: true}
		user, err := q.UpdateUser(ctx, req)

		if err != nil {
			logrus.Println("[Store UpdateUserTx] user is ", req.Email)
			logrus.Println("[Store UpdateUserTx] error create user is ", err)
			return err
		}

		logrus.Println("[Store UpdateUserTx] user ID is ", user.ID)
		arg := GetRoleUserByUserIDParams{
			UserID: req.ID,
			Limit:  1,
			Offset: 0,
		}
		getRoleID, err := q.GetRoleUserByUserID(ctx, arg)
		if err != nil {
			logrus.Println("[Store UpdateUserTx] error get role is ", err)
			return err
		}

		roleUserParams := UpdateRoleUserParams{
			ID:     getRoleID[0].ID,
			UserID: req.ID,
			RoleID: roleId,
			UpdatedBy: sql.NullInt64{
				Int64: userPayload.ID,
				Valid: true,
			},
		}

		roleUser, err := q.UpdateRoleUser(ctx, roleUserParams)
		if err != nil {
			logrus.Println("[Store UpdateUserTx] error update role is ", err)
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
