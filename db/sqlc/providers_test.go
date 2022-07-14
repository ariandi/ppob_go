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

func CreateRandomProvider(t *testing.T) Provider {
	user1 := CreateRandomUser(t)
	cstSh, _ := time.LoadLocation("Asia/Jakarta") //Jakarta
	layoutFormat := "2006-01-02 15:04:05"
	value := "2100-12-31 23:59:59"
	maxTime, _ := time.Parse(layoutFormat, value)

	arg := CreateProviderParams{
		User:      util.RandomUsername(),
		Name:      util.RandomUsername(),
		Secret:    util.RandomUsername() + strconv.FormatInt(util.RandomNumber(), 10),
		AddInfo1:  "Add info 1",
		AddInfo2:  "Add info 2",
		ValidFrom: sql.NullTime{Time: time.Now(), Valid: true},
		ValidTo:   sql.NullTime{Time: maxTime, Valid: true},
		BaseUrl: sql.NullString{
			String: "https://api.digiflazz.com/v1/",
			Valid:  true,
		},
		Method: sql.NullString{
			String: "POST",
			Valid:  true,
		},
		Inq: sql.NullString{
			String: "inq/transaction",
			Valid:  true,
		},
		Pay: sql.NullString{
			String: "pay/transaction",
			Valid:  true,
		},
		Adv: sql.NullString{
			String: "adv/transaction",
			Valid:  true,
		},
		Rev: sql.NullString{
			String: "rev/transaction",
			Valid:  true,
		},
		Status: "status",
		CreatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	provider, err := testQueries.CreateProvider(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, provider)

	require.Equal(t, arg.Name, provider.Name)
	require.Equal(t, arg.User, provider.User)
	require.Equal(t, arg.Secret, provider.Secret)
	require.Equal(t, arg.AddInfo1, provider.AddInfo1)
	require.Equal(t, arg.AddInfo2, provider.AddInfo2)
	require.Equal(t, arg.ValidFrom.Time.In(cstSh).Format("2006-01-02 15:04:05"), provider.ValidFrom.Time.In(cstSh).Format("2006-01-02 15:04:05"))
	require.Equal(t, arg.ValidTo.Time.Format("2006-01-02 15:04:05"), provider.ValidTo.Time.Format("2006-01-02 15:04:05"))
	require.Equal(t, arg.BaseUrl, provider.BaseUrl)
	require.Equal(t, arg.Method, provider.Method)
	require.Equal(t, arg.Inq, provider.Inq)
	require.Equal(t, arg.Pay, provider.Pay)
	//require.Equal(t, arg.Cmt, provider.Cmt)
	require.Equal(t, arg.Adv, provider.Adv)
	require.Equal(t, arg.Rev, provider.Rev)
	require.Equal(t, arg.Status, provider.Status)
	require.Equal(t, arg.CreatedBy.Int64, provider.CreatedBy.Int64)

	require.NotZero(t, provider.ID)
	//require.NotZero(t, category.Parent)
	require.NotZero(t, provider.CreatedAt)

	return provider
}

func TestCreateProvider(t *testing.T) {
	CreateRandomProvider(t)
}

func TestGetProvider(t *testing.T) {
	prov1 := CreateRandomProvider(t)
	prov2, err := testQueries.GetProvider(context.Background(), prov1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, prov2)

	require.Equal(t, prov2.Name, prov1.Name)
	require.Equal(t, prov2.User, prov1.User)
	require.Equal(t, prov2.Secret, prov1.Secret)
	require.Equal(t, prov2.AddInfo1, prov1.AddInfo1)
	require.Equal(t, prov2.AddInfo2, prov1.AddInfo2)
	require.Equal(t, prov2.ValidFrom, prov1.ValidFrom)
	require.Equal(t, prov2.ValidTo, prov1.ValidTo)
	require.Equal(t, prov2.BaseUrl, prov1.BaseUrl)
	require.Equal(t, prov2.Method, prov1.Method)
	require.Equal(t, prov2.Inq, prov1.Inq)
	require.Equal(t, prov2.Pay, prov1.Pay)
	//require.Equal(t, prov2.Cmt, prov1.Cmt)
	require.Equal(t, prov2.Adv, prov1.Adv)
	require.Equal(t, prov2.Rev, prov1.Rev)
	require.Equal(t, prov2.Status, prov1.Status)
	require.Equal(t, prov2.CreatedBy.Int64, prov1.CreatedBy.Int64)

	require.WithinDuration(
		t,
		prov1.CreatedAt.Time,
		prov2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateProvider(t *testing.T) {
	user1 := CreateRandomUser(t)
	prov1 := CreateRandomProvider(t)

	arg := UpdateProviderParams{
		ID:      prov1.ID,
		Name:    prov1.Name,
		SetName: true,
		SetCmt:  true,
		Cmt: sql.NullString{
			String: "Cmt",
			Valid:  true,
		},
		UpdatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	prov2, err := testQueries.UpdateProvider(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, prov2)
	require.NotEmpty(t, prov2.UpdatedAt.Time)
	require.NotEmpty(t, prov2.UpdatedBy.Int64)

	require.Equal(t, prov2.Name, prov1.Name)
	require.Equal(t, prov2.User, prov1.User)
	require.Equal(t, prov2.Secret, prov1.Secret)
	require.Equal(t, prov2.AddInfo1, prov1.AddInfo1)
	require.Equal(t, prov2.AddInfo2, prov1.AddInfo2)
	require.Equal(t, prov2.ValidFrom, prov1.ValidFrom)
	require.Equal(t, prov2.ValidTo, prov1.ValidTo)
	require.Equal(t, prov2.BaseUrl, prov1.BaseUrl)
	require.Equal(t, prov2.Method, prov1.Method)
	require.Equal(t, prov2.Inq, prov1.Inq)
	require.Equal(t, prov2.Pay, prov1.Pay)
	//require.Equal(t, prov2.Cmt, prov1.Cmt)
	require.Equal(t, prov2.Adv, prov1.Adv)
	require.Equal(t, prov2.Rev, prov1.Rev)
	require.Equal(t, prov2.Status, prov1.Status)
	require.Equal(t, prov2.CreatedBy.Int64, prov1.CreatedBy.Int64)

	require.WithinDuration(
		t,
		prov1.CreatedAt.Time,
		prov2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateDeleteProvider(t *testing.T) {
	user1 := CreateRandomUser(t)
	prov1 := CreateRandomProvider(t)

	var arg = UpdateInactiveProviderParams{
		ID:        prov1.ID,
		DeletedBy: sql.NullInt64{Int64: user1.ID, Valid: true},
	}
	prov2, err := testQueries.UpdateInactiveProvider(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, prov2)
	require.NotEmpty(t, prov2.DeletedAt.Time)
	require.NotEmpty(t, prov2.DeletedBy.Int64)

	require.Equal(t, prov1.Name, prov2.Name)
	require.Equal(t, prov1.User, prov2.User)
	require.Equal(t, prov1.CreatedBy.Int64, prov2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		prov1.CreatedAt.Time,
		prov2.CreatedAt.Time, time.Second,
	)
}

func TestDeleteProvider(t *testing.T) {
	prov1 := CreateRandomProvider(t)
	err := testQueries.DeleteProvider(context.Background(), prov1.ID)

	prov2, err := testQueries.GetProvider(context.Background(), prov1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, prov2)
}

func TestListProvider(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomProvider(t)
	}

	arg := ListProviderParams{
		Limit:  5,
		Offset: 5,
	}

	providers, err := testQueries.ListProvider(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, providers, 5)

	for _, provider := range providers {
		require.NotEmpty(t, provider)
	}
}
