package services

import (
	"encoding/json"
	"fmt"
	"github.com/adjust/rmq/v4"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/util"
	"strconv"
)

// UserService is
type UserService struct {
}

var userService *UserService
var redisConn rmq.Connection

// GetUserService is
func GetUserService(config util.Config) *UserService {

	if userService == nil {

		userService = new(UserService)

		redisDb, errRedisDb := strconv.Atoi(config.RedisDB)
		if errRedisDb != nil {
			fmt.Println("=======================================")
			fmt.Println("Cannot get redis db config : ", errRedisDb)
			fmt.Println("=======================================")
		}

		var err error
		redisConn, err = rmq.OpenConnection("redisService", "tcp", config.RedisUrl, redisDb, nil)

		if err != nil {
			fmt.Println("=======================================")
			fmt.Println("Error connect Redis : ", err)
			fmt.Println("=======================================")
		}
	}

	return userService
}

func (o *UserService) TestRedisMq(msg dto.LoginUserRequest) ([]string, error) {

	ret := []string{}
	queueName := "test_123456"
	redisQueue, err := redisConn.OpenQueue(queueName)
	if err != nil {
		return ret, err
	}

	//messageItem := msg
	byt, err := json.Marshal(msg)
	if err != nil {
		return ret, err
	}

	err = redisQueue.Publish(string(byt))
	ret = append(ret, "Success")

	return ret, nil
}

func (o *UserService) ValidateUserRole(user db.GetUserByUsernameRow) error {
	if user.RoleID.Int64 != 1 {
		return fmt.Errorf("you not allow to access this service")
	}

	return nil
}
