package db

import (
	"context"
	"database/sql"
	"github.com/ariandi/ppob_go/util"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomRole(t *testing.T, delete bool, softDelete bool) Role {
	var checkExistRole = false
	var role1 Role
	user1 := CreateRandomUser(t)
	roleName := util.RandomRole()

	// just for delete
	if delete {
		roleName = "delete"
	}

	if softDelete {
		roleName = "soft delete"
	}

	roles := GetRoles(t)

	for _, role2 := range roles {
		if role2.Name == roleName {
			checkExistRole = true
			role1 = role2
			break
		}
	}

	if checkExistRole {
		return role1
	}

	arg := CreateRoleParams{
		Name:  roleName,
		Level: 1,
		CreatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	role, err := testQueries.CreateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role)

	require.Equal(t, arg.Name, role.Name)
	require.Equal(t, arg.Level, role.Level)
	require.Equal(t, arg.CreatedBy.Int64, role.CreatedBy.Int64)

	require.NotZero(t, role.ID)
	require.NotZero(t, role.CreatedAt)

	return role
}

func TestCreateRole(t *testing.T) {
	CreateRandomRole(t, false, false)
}

func TestGetRole(t *testing.T) {
	role1 := CreateRandomRole(t, false, false)
	role2, err := testQueries.GetRole(context.Background(), role1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.Name, role2.Name)
	require.Equal(t, role1.Level, role2.Level)
	require.Equal(t, role1.CreatedBy.Int64, role2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		role1.CreatedAt.Time,
		role2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateRole(t *testing.T) {
	user1 := CreateRandomUser(t)
	role1 := CreateRandomRole(t, false, false)

	arg := UpdateRoleParams{
		ID:    role1.ID,
		Name:  role1.Name,
		Level: role1.Level,
		UpdatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	role2, err := testQueries.UpdateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role2)
	require.NotEmpty(t, role2.UpdatedAt.Time)
	require.NotEmpty(t, role2.UpdatedBy.Int64)

	require.Equal(t, role1.Name, role2.Name)
	require.Equal(t, role1.Level, role2.Level)
	require.Equal(t, role1.CreatedBy.Int64, role2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		role1.CreatedAt.Time,
		role2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateDeleteRole(t *testing.T) {
	user1 := CreateRandomUser(t)
	role1 := CreateRandomRole(t, false, true)

	var arg = UpdateInactiveRoleParams{
		ID:        role1.ID,
		DeletedBy: sql.NullInt64{Int64: user1.ID, Valid: true},
	}
	role2, err := testQueries.UpdateInactiveRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role2)
	require.NotEmpty(t, role2.DeletedAt.Time)
	require.NotEmpty(t, role2.DeletedBy.Int64)

	require.Equal(t, role1.Name, role2.Name)
	require.Equal(t, role1.Level, role2.Level)
	require.Equal(t, role1.CreatedBy.Int64, role2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		role1.CreatedAt.Time,
		role2.CreatedAt.Time, time.Second,
	)
}

func TestDeleteRole(t *testing.T) {
	role1 := CreateRandomRole(t, true, false)
	err := testQueries.DeleteRole(context.Background(), role1.ID)

	role2, err := testQueries.GetRole(context.Background(), role1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, role2)
}

func TestListRoles(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomRole(t, false, false)
	}

	arg := ListRoleParams{
		Limit:  5,
		Offset: 0,
	}

	roles, err := testQueries.ListRole(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, roles, 5)

	for _, role := range roles {
		require.NotEmpty(t, role)
	}
}

func GetRoles(t *testing.T) []Role {
	arg := ListRoleWithDeleteParams{
		Limit:  5,
		Offset: 0,
	}

	roles, err := testQueries.ListRoleWithDelete(context.Background(), arg)
	require.NoError(t, err)

	return roles
}
