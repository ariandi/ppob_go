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

func CreateRandomTransaction(t *testing.T) Transaction {
	user1 := CreateRandomUser(t)
	cat1 := CreateRandomCategory(t)
	prod1 := CreateRandomProduct(t)
	partner1 := CreateRandomPartner(t)
	prov1 := CreateRandomProvider(t)

	arg := CreateTransactionParams{
		TxID:       util.SetTxID(),
		BillID:     util.RandomNumberStr(),
		CustName:   sql.NullString{String: user1.Name, Valid: true},
		Amount:     sql.NullString{String: strconv.FormatFloat(100000.00, 'f', 2, 64), Valid: true},
		Admin:      sql.NullString{String: strconv.FormatFloat(100.00, 'f', 2, 64), Valid: true},
		TotAmount:  sql.NullString{String: strconv.FormatFloat(100100.00, 'f', 2, 64), Valid: true},
		FeePartner: sql.NullString{String: strconv.FormatFloat(100.00, 'f', 2, 64), Valid: true},
		FeePpob:    sql.NullString{String: strconv.FormatFloat(100.00, 'f', 2, 64), Valid: true},
		CatID: sql.NullInt64{
			Int64: cat1.ID,
			Valid: true,
		},
		CatName: sql.NullString{
			String: cat1.Name,
			Valid:  true,
		},
		ProdID: sql.NullInt64{
			Int64: prod1.ID,
			Valid: true,
		},
		ProdName: sql.NullString{
			String: prod1.Name,
			Valid:  true,
		},
		PartnerID: sql.NullInt64{
			Int64: partner1.ID,
			Valid: true,
		},
		PartnerName: sql.NullString{
			String: partner1.Name,
			Valid:  true,
		},
		ProviderID: sql.NullInt64{
			Int64: prov1.ID,
			Valid: true,
		},
		ProviderName: sql.NullString{
			String: prov1.Name,
			Valid:  true,
		},
		Status: util.RandomTrxStatus(),
		ReqInqParams: sql.NullString{
			String: "{}",
			Valid:  true,
		},
		CreatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
		RefID: "1234567890",
		FirstBalance: sql.NullString{
			String: "1000",
			Valid:  true,
		},
		LastBalance: sql.NullString{
			String: "1000",
			Valid:  true,
		},
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.TxID, transaction.TxID)
	require.Equal(t, arg.Status, transaction.Status)
	require.Equal(t, arg.BillID, transaction.BillID)
	require.Equal(t, arg.CustName, transaction.CustName)
	require.Equal(t, arg.Amount, transaction.Amount)
	require.Equal(t, arg.Admin, transaction.Admin)
	require.Equal(t, arg.TotAmount, transaction.TotAmount)
	require.Equal(t, arg.FeePartner, transaction.FeePartner)
	require.Equal(t, arg.FeePpob, transaction.FeePpob)
	require.Equal(t, arg.CatID, transaction.CatID)
	require.Equal(t, arg.CatName, transaction.CatName)
	require.Equal(t, arg.ProdID, transaction.ProdID)
	require.Equal(t, arg.ProdName, transaction.ProdName)
	require.Equal(t, arg.PartnerID, transaction.PartnerID)
	require.Equal(t, arg.PartnerName, transaction.PartnerName)
	require.Equal(t, arg.ProviderID, transaction.ProviderID)
	require.Equal(t, arg.ProviderName, transaction.ProviderName)
	require.Equal(t, arg.ReqInqParams, transaction.ReqInqParams)
	require.Equal(t, arg.CreatedBy.Int64, transaction.CreatedBy.Int64)

	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)

	return transaction
}

func TestCreateTransaction(t *testing.T) {
	CreateRandomTransaction(t)
}

func TestGetTransaction(t *testing.T) {
	trx1 := CreateRandomTransaction(t)
	trx2, err := testQueries.GetTransaction(context.Background(), trx1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trx2)

	require.Equal(t, trx1.TxID, trx2.TxID)
	require.Equal(t, trx1.Status, trx2.Status)
	require.Equal(t, trx1.BillID, trx2.BillID)
	require.Equal(t, trx1.CustName, trx2.CustName)
	require.Equal(t, trx1.Amount, trx2.Amount)
	require.Equal(t, trx1.Admin, trx2.Admin)
	require.Equal(t, trx1.TotAmount, trx2.TotAmount)
	require.Equal(t, trx1.FeePartner, trx2.FeePartner)
	require.Equal(t, trx1.FeePpob, trx2.FeePpob)
	require.Equal(t, trx1.CatID, trx2.CatID)
	require.Equal(t, trx1.CatName, trx2.CatName)
	require.Equal(t, trx1.ProdID, trx2.ProdID)
	require.Equal(t, trx1.ProdName, trx2.ProdName)
	require.Equal(t, trx1.PartnerID, trx2.PartnerID)
	require.Equal(t, trx1.PartnerName, trx2.PartnerName)
	require.Equal(t, trx1.ProviderID, trx2.ProviderID)
	require.Equal(t, trx1.ProviderName, trx2.ProviderName)
	require.Equal(t, trx1.ReqInqParams, trx2.ReqInqParams)
	require.Equal(t, trx1.CreatedBy.Int64, trx2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		trx1.CreatedAt.Time,
		trx2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateTransaction(t *testing.T) {
	user1 := CreateRandomUser(t)
	trx1 := CreateRandomTransaction(t)

	arg := UpdateTransactionParams{
		ID:     trx1.ID,
		Status: trx1.Status,
		ResInqParams: sql.NullString{
			String: "{}",
			Valid:  true,
		},
		ReqPayParams: sql.NullString{
			String: "{}",
			Valid:  true,
		},
		ResPayParams: sql.NullString{
			String: "{}",
			Valid:  true,
		},
		SetStatus:       true,
		SetResInqParams: true,
		SetReqPayParams: true,
		SetResPayParams: true,
		UpdatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	trx2, err := testQueries.UpdateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trx2)
	require.NotEmpty(t, trx2.UpdatedAt.Time)
	require.NotEmpty(t, trx2.UpdatedBy.Int64)

	require.Equal(t, trx1.TxID, trx2.TxID)
	require.Equal(t, trx1.Status, trx2.Status)
	require.Equal(t, arg.ResInqParams.String, trx2.ResInqParams.String)
	require.Equal(t, arg.ReqPayParams.String, trx2.ReqPayParams.String)
	require.Equal(t, arg.ResPayParams.String, trx2.ResPayParams.String)
	require.Equal(t, trx1.CreatedBy.Int64, trx2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		trx1.CreatedAt.Time,
		trx2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateDeleteTransaction(t *testing.T) {
	user1 := CreateRandomUser(t)
	trx1 := CreateRandomTransaction(t)

	var arg = UpdateInactiveTransactionParams{
		ID:        trx1.ID,
		DeletedBy: sql.NullInt64{Int64: user1.ID, Valid: true},
	}
	trx2, err := testQueries.UpdateInactiveTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trx2)
	require.NotEmpty(t, trx2.DeletedAt.Time)
	require.NotEmpty(t, trx2.DeletedBy.Int64)

	require.Equal(t, trx1.Status, trx2.Status)
	require.Equal(t, trx1.CreatedBy.Int64, trx2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		trx1.CreatedAt.Time,
		trx2.CreatedAt.Time, time.Second,
	)
}

func TestDeleteTransaction(t *testing.T) {
	trx1 := CreateRandomTransaction(t)
	err := testQueries.DeleteTransaction(context.Background(), trx1.ID)

	trx2, err := testQueries.GetTransaction(context.Background(), trx1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, trx2)
}

func TestListTransaction(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransaction(t)
	}

	arg := ListTransactionParams{
		Limit:  5,
		Offset: 5,
	}

	trans, err := testQueries.ListTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, trans, 5)

	for _, trx := range trans {
		require.NotEmpty(t, trx)
	}
}
