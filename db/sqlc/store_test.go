package db

import (
	"context"
	"database/sql"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateUserTx(t *testing.T) {
	store := NewStore(testDB)
	user := CreateRandomUser(t)
	n := 5

	errs := make(chan error)
	results := make(chan dto.UserResponse)
	users := make(chan dto.CreateUserRequest)
	for i := 0; i < n; i++ {
		nameTmp := util.RandomString(6 + i)
		userTmp := dto.CreateUserRequest{
			Name:           nameTmp,
			Email:          nameTmp + "@gmail.com",
			Username:       nameTmp,
			Password:       nameTmp,
			Phone:          "081283743874",
			IdentityNumber: "1234123412341234",
			RoleID:         int64(1),
		}
		go func() {
			users <- userTmp
			//fmt.Println(">> user email is :", userTmp.Email)
			//fmt.Println(">> user is :", userTmp)
			arg := CreateUserParams{
				Name:     userTmp.Name,
				Email:    userTmp.Email,
				Username: userTmp.Username,
				Password: sql.NullString{
					String: userTmp.Password,
					Valid:  true,
				},
				Phone:          userTmp.Phone,
				IdentityNumber: userTmp.IdentityNumber,
			}

			payload := new(token.Payload)
			payload.ID = uuid.New()
			payload.Username = user.Username
			payload.UserID = user.ID
			payload.ExpiredAt = time.Now()
			payload.IssuedAt = time.Now()

			result, err := store.CreateUserTx(context.Background(), arg, payload, userTmp.RoleID)

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		user2 := <-users
		require.NotEmpty(t, user2)

		err := <-errs
		require.NoError(t, err)

		userResult := <-results
		require.NotEmpty(t, userResult)

		require.Equal(t, user2.Username, userResult.Username)
		require.Equal(t, user2.Name, userResult.Name)
		require.Equal(t, user2.Email, userResult.Email)
		require.NotZero(t, userResult.ID)
		require.NotEmpty(t, userResult.Role)

		_, err = store.GetUser(context.Background(), userResult.ID)
		require.NoError(t, err)
	}
}

//func TestGoRoutin(t *testing.T) {
//	n := 5
//	aduh := make(chan string)
//	aduh2 := make(chan string)
//	aduh3 := make(chan string)
//	for i := 0; i < n; i++ {
//		go func() {
//			aduh <- "test"
//			aduh2 <- "test2"
//			aduh3 <- "test3"
//		}()
//	}
//
//	aduhResult := <-aduh
//	aduhResult2 := <-aduh2
//	aduhResult3 := <-aduh3
//	fmt.Println("wewe", aduhResult)
//	fmt.Println("wewe", aduhResult2)
//	fmt.Println("wewe", aduhResult3)
//	require.NotEmpty(t, aduhResult)
//
//}
