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

func CreateRandomCategory(t *testing.T) Category {
	user1 := CreateRandomUser(t)

	arg := CreateCategoryParams{
		Name: util.RandomUsername(),
		UpSelling: sql.NullString{
			String: "100",
			Valid:  true,
		},
		CreatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	category, err := testQueries.CreateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.CreatedBy.Int64, category.CreatedBy.Int64)

	require.NotZero(t, category.ID)
	//require.NotZero(t, category.Parent)
	require.NotZero(t, category.CreatedAt)

	return category
}

func TestCreateCategory(t *testing.T) {
	CreateRandomCategory(t)
}

func TestGetCategory(t *testing.T) {
	cat1 := CreateRandomCategory(t)
	cat2, err := testQueries.GetCategory(context.Background(), cat1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, cat2)

	require.Equal(t, cat1.Name, cat2.Name)
	require.Equal(t, cat1.Parent, cat2.Parent)
	require.Equal(t, cat1.CreatedBy.Int64, cat2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		cat1.CreatedAt.Time,
		cat2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateCategory(t *testing.T) {
	user1 := CreateRandomUser(t)
	cat1 := CreateRandomCategory(t)

	arg := UpdateCategoryParams{
		ID:        cat1.ID,
		Name:      cat1.Name,
		SetName:   true,
		SetParent: false,
		UpdatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	cat2, err := testQueries.UpdateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cat2)
	require.NotEmpty(t, cat2.UpdatedAt.Time)
	require.NotEmpty(t, cat2.UpdatedBy.Int64)

	require.Equal(t, cat1.Name, cat2.Name)
	require.Equal(t, cat1.Parent, cat2.Parent)
	require.Equal(t, cat1.CreatedBy.Int64, cat2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		cat1.CreatedAt.Time,
		cat2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateDeleteCategory(t *testing.T) {
	user1 := CreateRandomUser(t)
	cat1 := CreateRandomCategory(t)

	var arg = UpdateInactiveCategoryParams{
		ID:        cat1.ID,
		DeletedBy: sql.NullInt64{Int64: user1.ID, Valid: true},
	}
	cat2, err := testQueries.UpdateInactiveCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cat2)
	require.NotEmpty(t, cat2.DeletedAt.Time)
	require.NotEmpty(t, cat2.DeletedBy.Int64)

	require.Equal(t, cat1.Name, cat2.Name)
	require.Equal(t, cat1.Parent, cat2.Parent)
	require.Equal(t, cat1.CreatedBy.Int64, cat2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		cat1.CreatedAt.Time,
		cat2.CreatedAt.Time, time.Second,
	)
}

func TestDeleteCategory(t *testing.T) {
	cat1 := CreateRandomCategory(t)
	err := testQueries.DeleteCategories(context.Background(), cat1.ID)

	cat2, err := testQueries.GetCategory(context.Background(), cat1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, cat2)
}

func TestListCategories(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomCategory(t)
	}

	arg := ListCategoryParams{
		Limit:  5,
		Offset: 5,
	}

	categories, err := testQueries.ListCategory(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, categories, 5)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}
