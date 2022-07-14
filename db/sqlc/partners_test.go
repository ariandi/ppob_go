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

func CreateRandomPartner(t *testing.T) Partner {
	user1 := CreateRandomUser(t)
	cstSh, _ := time.LoadLocation("Asia/Jakarta") //Jakarta
	layoutFormat := "2006-01-02 15:04:05"
	value := "2100-12-31 23:59:59"
	maxTime, _ := time.Parse(layoutFormat, value)

	arg := CreatePartnerParams{
		User:        util.RandomUsername(),
		Name:        util.RandomUsername(),
		Secret:      util.RandomUsername() + strconv.FormatInt(util.RandomNumber(), 10),
		AddInfo1:    "Add info 1",
		AddInfo2:    "Add info 2",
		ValidFrom:   sql.NullTime{Time: time.Now(), Valid: true},
		ValidTo:     sql.NullTime{Time: maxTime, Valid: true},
		PaymentType: "deposit",
		Status:      "status",
		CreatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	partner, err := testQueries.CreatePartner(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, partner)

	require.Equal(t, arg.Name, partner.Name)
	require.Equal(t, arg.User, partner.User)
	require.Equal(t, arg.Secret, partner.Secret)
	require.Equal(t, arg.AddInfo1, partner.AddInfo1)
	require.Equal(t, arg.AddInfo2, partner.AddInfo2)
	require.Equal(t, arg.ValidFrom.Time.In(cstSh).Format("2006-01-02 15:04:05"), partner.ValidFrom.Time.In(cstSh).Format("2006-01-02 15:04:05"))
	require.Equal(t, arg.ValidTo.Time.Format("2006-01-02 15:04:05"), partner.ValidTo.Time.Format("2006-01-02 15:04:05"))
	require.Equal(t, arg.PaymentType, partner.PaymentType)
	require.Equal(t, arg.Status, partner.Status)
	require.Equal(t, arg.CreatedBy.Int64, partner.CreatedBy.Int64)

	require.NotZero(t, partner.ID)
	//require.NotZero(t, category.Parent)
	require.NotZero(t, partner.CreatedAt)

	return partner
}

func TestCreatePartner(t *testing.T) {
	CreateRandomPartner(t)
}

func TestGetPartner(t *testing.T) {
	partner1 := CreateRandomPartner(t)
	partner2, err := testQueries.GetPartner(context.Background(), partner1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, partner2)

	require.Equal(t, partner2.Name, partner1.Name)
	require.Equal(t, partner2.User, partner1.User)
	require.Equal(t, partner2.Secret, partner1.Secret)
	require.Equal(t, partner2.AddInfo1, partner1.AddInfo1)
	require.Equal(t, partner2.AddInfo2, partner1.AddInfo2)
	require.Equal(t, partner2.ValidFrom, partner1.ValidFrom)
	require.Equal(t, partner2.ValidTo, partner1.ValidTo)
	require.Equal(t, partner2.PaymentType, partner1.PaymentType)
	require.Equal(t, partner2.Status, partner1.Status)
	require.Equal(t, partner2.CreatedBy.Int64, partner1.CreatedBy.Int64)

	require.WithinDuration(
		t,
		partner1.CreatedAt.Time,
		partner2.CreatedAt.Time, time.Second,
	)
}

func TestUpdatePartner(t *testing.T) {
	user1 := CreateRandomUser(t)
	partner1 := CreateRandomPartner(t)

	arg := UpdatePartnerParams{
		ID:      partner1.ID,
		Name:    partner1.Name,
		SetName: true,
		UpdatedBy: sql.NullInt64{
			Int64: user1.ID,
			Valid: true,
		},
	}

	partner2, err := testQueries.UpdatePartner(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, partner2)
	require.NotEmpty(t, partner2.UpdatedAt.Time)
	require.NotEmpty(t, partner2.UpdatedBy.Int64)

	require.Equal(t, partner1.Name, partner2.Name)
	require.Equal(t, partner1.User, partner2.User)
	require.Equal(t, partner1.Secret, partner2.Secret)
	require.Equal(t, partner1.AddInfo1, partner2.AddInfo1)
	require.Equal(t, partner1.AddInfo2, partner2.AddInfo2)
	require.Equal(t, partner1.ValidFrom, partner2.ValidFrom)
	require.Equal(t, partner1.ValidTo, partner2.ValidTo)
	require.Equal(t, partner1.PaymentType, partner2.PaymentType)
	require.Equal(t, partner1.Status, partner2.Status)
	require.Equal(t, partner1.CreatedBy.Int64, partner2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		partner1.CreatedAt.Time,
		partner2.CreatedAt.Time, time.Second,
	)
}

func TestUpdateDeletePartner(t *testing.T) {
	user1 := CreateRandomUser(t)
	partner1 := CreateRandomPartner(t)

	var arg = UpdateInactivePartnerParams{
		ID:        partner1.ID,
		DeletedBy: sql.NullInt64{Int64: user1.ID, Valid: true},
	}
	partner2, err := testQueries.UpdateInactivePartner(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, partner2)
	require.NotEmpty(t, partner2.DeletedAt.Time)
	require.NotEmpty(t, partner2.DeletedBy.Int64)

	require.Equal(t, partner1.Name, partner2.Name)
	require.Equal(t, partner1.User, partner2.User)
	require.Equal(t, partner1.CreatedBy.Int64, partner2.CreatedBy.Int64)

	require.WithinDuration(
		t,
		partner1.CreatedAt.Time,
		partner2.CreatedAt.Time, time.Second,
	)
}

func TestDeletePartner(t *testing.T) {
	partner1 := CreateRandomPartner(t)
	err := testQueries.DeletePartner(context.Background(), partner1.ID)

	partner2, err := testQueries.GetPartner(context.Background(), partner1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, partner2)
}

func TestListPartner(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomPartner(t)
	}

	arg := ListPartnerParams{
		Limit:  5,
		Offset: 5,
	}

	partners, err := testQueries.ListPartner(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, partners, 5)

	for _, partner := range partners {
		require.NotEmpty(t, partner)
	}
}
