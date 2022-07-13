package db

import (
	"context"
	"database/sql"
	"github.com/ariandi/ppob_go/util"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func CreateRandomUser(t *testing.T) User {
	username := util.RandomUsername()

	arg := CreateUserParams{
		Name:      username,
		Email:     username + "@gmial.com",
		Password:  sql.NullString{String: strconv.FormatInt(util.RandomNumber(), 10), Valid: true},
		Username:  username,
		CreatedBy: sql.NullInt64{Int64: 1, Valid: true},
		Phone:     strconv.FormatInt(util.RandomNumber(), 10),
		Balance: sql.NullString{
			String: "0.00",
			Valid:  true,
		},
		IdentityNumber: strconv.FormatInt(util.RandomNumber(), 10),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password.String, user.Password.String)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Phone, user.Phone)
	require.Equal(t, arg.CreatedBy.Int64, user.CreatedBy.Int64)
	require.Equal(t, arg.Balance.String, user.Balance.String)
	require.Equal(t, arg.IdentityNumber, user.IdentityNumber)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password.String, user2.Password.String)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.CreatedBy.Int64, user2.CreatedBy.Int64)
	require.Equal(t, user1.Balance.String, user2.Balance.String)
	require.Equal(t, user1.IdentityNumber, user2.IdentityNumber)

	require.WithinDuration(
		t,
		user1.CreatedAt.Time,
		user2.CreatedAt.Time, time.Second,
	)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password.String, user2.Password.String)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.CreatedBy.Int64, user2.CreatedBy.Int64)
	require.Equal(t, user1.Balance.String, user2.Balance.String)
	require.Equal(t, user1.IdentityNumber, user2.IdentityNumber)

	require.WithinDuration(
		t,
		user1.CreatedAt.Time,
		user2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateUser(t *testing.T) {
	user1 := CreateRandomUser(t)

	arg := UpdateUserParams{
		ID:          user1.ID,
		SetName:     true,
		Name:        user1.Name,
		SetPassword: true,
		Password: sql.NullString{
			String: user1.Password.String,
			Valid:  true,
		},
		UpdatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.NotEmpty(t, user2.UpdatedAt.Time)
	require.NotEmpty(t, user2.UpdatedBy.Int64)

	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password.String, user2.Password.String)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.CreatedBy.Int64, user2.CreatedBy.Int64)
	require.Equal(t, user1.Balance.String, user2.Balance.String)
	require.Equal(t, user1.IdentityNumber, user2.IdentityNumber)

	require.WithinDuration(
		t,
		user1.CreatedAt.Time,
		user2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateDeleteUser(t *testing.T) {
	user1 := CreateRandomUser(t)

	arg := UpdateInactiveUserParams{
		ID:        user1.ID,
		DeletedBy: sql.NullInt64{Int64: user1.ID, Valid: true},
	}

	user2, err := testQueries.UpdateInactiveUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.NotEmpty(t, user2.DeletedAt.Time)
	require.NotEmpty(t, user2.DeletedBy.Int64)

	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password.String, user2.Password.String)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.Equal(t, user1.CreatedBy.Int64, user2.CreatedBy.Int64)
	require.Equal(t, user1.Balance.String, user2.Balance.String)
	require.Equal(t, user1.IdentityNumber, user2.IdentityNumber)

	require.WithinDuration(
		t,
		user1.CreatedAt.Time,
		user2.CreatedAt.Time, time.Second,
	)
}

func TestDeleteUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomUser(t)
	}

	arg := ListUserParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUser(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
