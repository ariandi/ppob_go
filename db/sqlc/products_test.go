package db

import (
	"context"
	"database/sql"
	"github.com/ariandi/ppob_go/util"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"math"
	"strconv"
	"testing"
	"time"
)

func CreateRandomProduct(t *testing.T) Product {
	user1 := CreateRandomUser(t)
	cat1 := CreateRandomCategory(t)
	prov1 := CreateRandomProvider(t)
	amountArg := 100.05

	arg := CreateProductParams{
		CatID:      cat1.ID,
		Status:     util.RandomStatus(),
		Name:       util.RandomUsername(),
		Parent:     0,
		Amount:     strconv.FormatFloat(amountArg, 'f', 6, 64),
		ProviderID: prov1.ID,
		CreatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	amount, err := strconv.ParseFloat(product.Amount, 64)
	require.NoError(t, err)

	amount = math.Floor(amount*100) / 100 // two digit behind comma
	amountStr := strconv.FormatFloat(amount, 'f', 6, 64)

	require.Equal(t, arg.Name, product.Name)
	require.Equal(t, arg.CatID, product.CatID)
	require.Equal(t, arg.ProviderID, product.ProviderID)
	require.Equal(t, arg.Status, product.Status)
	require.Equal(t, arg.Amount, amountStr)
	require.Equal(t, arg.Parent, product.Parent)
	require.Equal(t, arg.CreatedBy.Int64, product.CreatedBy.Int64)

	require.NotZero(t, product.ID)
	//require.NotZero(t, category.Parent)
	require.NotZero(t, product.CreatedAt)

	return product
}

func TestCreateProduct(t *testing.T) {
	CreateRandomProduct(t)
}

func TestGetProduct(t *testing.T) {
	prod1 := CreateRandomProduct(t)
	prod2, err := testQueries.GetProduct(context.Background(), prod1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, prod2)

	require.Equal(t, prod1.Name, prod2.Name)
	require.Equal(t, prod1.CatID, prod2.CatID)
	require.Equal(t, prod1.ProviderID, prod2.ProviderID)
	require.Equal(t, prod1.Status, prod2.Status)
	require.Equal(t, prod1.Amount, prod2.Amount)
	require.Equal(t, prod1.Parent, prod2.Parent)
	require.Equal(t, prod1.CreatedBy.Int64, prod2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		prod1.CreatedAt.Time,
		prod2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateProduct(t *testing.T) {
	user1 := CreateRandomUser(t)
	cat1 := CreateRandomCategory(t)
	prov1 := CreateRandomProvider(t)
	prod1 := CreateRandomProduct(t)

	arg := UpdateProductParams{
		ID:         prod1.ID,
		CatID:      cat1.ID,
		ProviderID: prov1.ID,
		Name:       prod1.Name,
		//Status:     prod1.Status,
		Amount:      strconv.Itoa(200.00),
		SetName:     true,
		SetCat:      true,
		SetProvider: true,
		SetStatus:   false,
		SetAmount:   false,
		SetParent:   false,
		UpdatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	prod2, err := testQueries.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, prod2)
	require.NotEmpty(t, prod2.UpdatedAt.Time)
	require.NotEmpty(t, prod2.UpdatedBy.Int64)

	require.Equal(t, prod1.Name, prod2.Name)
	require.Equal(t, cat1.ID, prod2.CatID)
	require.Equal(t, prov1.ID, prod2.ProviderID)
	require.Equal(t, prod1.Status, prod2.Status)
	require.Equal(t, prod1.Amount, prod2.Amount)
	require.Equal(t, prod1.Parent, prod2.Parent)
	require.Equal(t, prod1.CreatedBy.Int64, prod2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		cat1.CreatedAt.Time,
		prod2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateDeleteProduct(t *testing.T) {
	user1 := CreateRandomUser(t)
	prod1 := CreateRandomProduct(t)

	var arg = UpdateInactiveProductParams{
		ID:        prod1.ID,
		DeletedBy: sql.NullInt64{Int64: user1.ID, Valid: true},
	}
	prod2, err := testQueries.UpdateInactiveProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, prod2)
	require.NotEmpty(t, prod2.DeletedAt.Time)
	require.NotEmpty(t, prod2.DeletedBy.Int64)

	require.Equal(t, prod1.Name, prod2.Name)
	require.Equal(t, prod1.CatID, prod2.CatID)
	require.Equal(t, prod1.ProviderID, prod2.ProviderID)
	require.Equal(t, prod1.Status, prod2.Status)
	require.Equal(t, prod1.Amount, prod2.Amount)
	require.Equal(t, prod1.Parent, prod2.Parent)
	require.Equal(t, prod1.CreatedBy.Int64, prod2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		prod1.CreatedAt.Time,
		prod2.CreatedAt.Time, time.Second,
	)
}

func TestDeleteProduct(t *testing.T) {
	prod1 := CreateRandomProduct(t)
	err := testQueries.DeleteProduct(context.Background(), prod1.ID)

	prod2, err := testQueries.GetProduct(context.Background(), prod1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, prod2)
}

func TestListProduct(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomProduct(t)
	}

	arg := ListProductParams{
		Limit:  5,
		Offset: 5,
	}

	products, err := testQueries.ListProduct(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, products, 5)

	for _, product := range products {
		require.NotEmpty(t, product)
	}
}
