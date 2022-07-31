package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomRoleUser(t *testing.T) RoleUser {
	user1 := CreateRandomUser(t)
	role1 := CreateRandomRole(t, false, false)

	arg := CreateRoleUserParams{
		UserID: user1.ID,
		RoleID: role1.ID,
		CreatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	roleUser, err := testQueries.CreateRoleUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, roleUser)

	require.Equal(t, arg.UserID, roleUser.UserID)
	require.Equal(t, arg.RoleID, roleUser.RoleID)
	require.Equal(t, arg.CreatedBy.Int64, roleUser.CreatedBy.Int64)

	require.NotZero(t, roleUser.ID)
	require.NotZero(t, roleUser.CreatedAt)

	return roleUser
}

func TestCreateRoleUser(t *testing.T) {
	CreateRandomRoleUser(t)
}

func TestGetRoleUser(t *testing.T) {
	roleUser1 := CreateRandomRoleUser(t)
	roleUser2, err := testQueries.GetRoleUserByID(context.Background(), roleUser1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, roleUser2)

	require.Equal(t, roleUser1.UserID, roleUser2.UserID)
	require.Equal(t, roleUser1.RoleID, roleUser2.RoleID)
	require.Equal(t, roleUser1.CreatedBy.Int64, roleUser2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		roleUser1.CreatedAt.Time,
		roleUser2.CreatedAt.Time, time.Second,
	)
}

func TestGetRoleUserByUserID(t *testing.T) {
	roleUser1 := CreateRandomRoleUser(t)

	arg := GetRoleUserByUserIDParams{
		UserID: roleUser1.UserID,
		Limit:  5,
		Offset: 0,
	}

	roleUsers, err := testQueries.GetRoleUserByUserID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, roleUsers)

	for _, roleUser := range roleUsers {
		require.NotEmpty(t, roleUser)
	}
}

func TestGetRoleUserByRoleID(t *testing.T) {
	roleUser1 := CreateRandomRoleUser(t)

	arg := GetRoleUserByRoleIDParams{
		RoleID: roleUser1.RoleID,
		Limit:  5,
		Offset: 0,
	}

	roleUsers, err := testQueries.GetRoleUserByRoleID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, roleUsers)

	for _, roleUser := range roleUsers {
		require.NotEmpty(t, roleUser)
	}
}

func TestUpdateRoleUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	role1 := CreateRandomRole(t, false, false)
	roleUser1 := CreateRandomRoleUser(t)

	arg := UpdateRoleUserParams{
		ID:     roleUser1.ID,
		UserID: user1.ID,
		RoleID: role1.ID,
		UpdatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	roleUser2, err := testQueries.UpdateRoleUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, roleUser2)
	require.NotEmpty(t, roleUser2.UpdatedAt.Time)
	require.NotEmpty(t, roleUser2.UpdatedBy.Int64)

	require.Equal(t, user1.ID, roleUser2.UserID)
	require.Equal(t, role1.ID, roleUser2.RoleID)
	require.Equal(t, user1.ID, roleUser2.UpdatedBy.Int64)

	require.WithinDuration(
		t,
		roleUser1.CreatedAt.Time,
		roleUser2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateDeleteRoleUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	roleUser1 := CreateRandomRoleUser(t)

	var arg = UpdateInactiveRoleUserParams{
		ID:        roleUser1.ID,
		DeletedBy: sql.NullInt64{Int64: user1.ID, Valid: true},
	}
	roleUser2, err := testQueries.UpdateInactiveRoleUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, roleUser2)
	require.NotEmpty(t, roleUser2.DeletedAt.Time)
	require.NotEmpty(t, roleUser2.DeletedBy.Int64)

	require.Equal(t, roleUser1.UserID, roleUser2.UserID)
	require.Equal(t, roleUser1.RoleID, roleUser2.RoleID)
	require.Equal(t, roleUser1.CreatedBy.Int64, roleUser2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		roleUser1.CreatedAt.Time,
		roleUser2.CreatedAt.Time, time.Second,
	)
}

func TestDeleteRoleUser(t *testing.T) {
	roleUser1 := CreateRandomRoleUser(t)
	err := testQueries.DeleteRoleUser(context.Background(), roleUser1.ID)

	roleUser2, err := testQueries.GetRoleUserByID(context.Background(), roleUser1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, roleUser2)
}

func TestListRolesUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomRoleUser(t)
	}

	arg := ListRoleUserParams{
		Limit:  5,
		Offset: 0,
	}

	roleUsers, err := testQueries.ListRoleUser(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, roleUsers, 5)

	for _, roleUser := range roleUsers {
		require.NotEmpty(t, roleUser)
	}
}
