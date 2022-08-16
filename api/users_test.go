package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/ariandi/ppob_go/db/mock"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.Password.String)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)
	userResponse := randomUserResponse(user)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":            user.Name,
				"email":           user.Email,
				"username":        user.Username,
				"phone":           user.Phone,
				"identity_number": user.IdentityNumber,
				"password":        password,
				"role_id":         user.RoleID.Int64,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Name:           user.Name,
					Email:          user.Email,
					Username:       user.Username,
					Phone:          user.Phone,
					IdentityNumber: user.IdentityNumber,
					//BankCode: sql.NullInt64{
					//	Int64: 1,
					//	Valid: true,
					//},
					Balance: sql.NullString{
						String: "0.00",
						Valid:  true,
					},
					CreatedBy: sql.NullInt64{
						Int64: 0,
						Valid: true,
					},
				}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserParams(arg, password), gomock.Any(), user.RoleID.Int64).
					Times(1).
					Return(userResponse, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, userResponse)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"name":            user.Name,
				"email":           user.Email,
				"username":        user.Username,
				"phone":           user.Phone,
				"identity_number": user.IdentityNumber,
				"password":        password,
				"role_id":         1,
				//"created_by":      user.CreatedBy,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				//store.EXPECT().
				//	GetUserByUsername(gomock.Any(), gomock.Any()).
				//	Times(0).
				//	Return(user, nil)
				//
				//store.EXPECT().
				//	CreateUser(gomock.Any(), gomock.Any()).
				//	Times(0).
				//	Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		//{
		//	name: "UnauthorizedUser",
		//	body: gin.H{
		//		"name":            user.Name,
		//		"email":           user.Email,
		//		"username":        user.Username,
		//		"phone":           user.Phone,
		//		"identity_number": user.IdentityNumber,
		//		"password":        password,
		//		"role_id":         1,
		//	},
		//	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		//		addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, "unauthorized_user", time.Minute)
		//	},
		//	buildStubs: func(store *mockdb.MockStore) {
		//		store.EXPECT().
		//			CreateUserTx(gomock.Any(), gomock.Eq(user.Username), gomock.Any()).
		//			Times(1).
		//			Return(user, nil)
		//		store.EXPECT().
		//			CreateUser(gomock.Any(), gomock.Any()).
		//			Times(1).
		//			Return(user, nil)
		//	},
		//	checkResponse: func(recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		//	},
		//},
		{
			name: "InternalError",
			body: gin.H{
				"name":            user.Name,
				"email":           user.Email,
				"username":        user.Username,
				"phone":           user.Phone,
				"identity_number": user.IdentityNumber,
				"password":        password,
				"role_id":         user.RoleID.Int64,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any(), gomock.Any(), user.RoleID.Int64).
					Times(1).
					Return(dto.UserResponse{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"name":            user.Name,
				"email":           user.Email,
				"username":        user.Username,
				"phone":           user.Phone,
				"identity_number": user.IdentityNumber,
				"password":        password,
				"role_id":         user.RoleID.Int64,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any(), gomock.Any(), user.RoleID.Int64).
					Times(1).
					Return(dto.UserResponse{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"name":            user.Name,
				"email":           user.Email,
				"username":        "invalid-user#1",
				"phone":           user.Phone,
				"identity_number": user.IdentityNumber,
				"password":        password,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0).
					Return(user, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"name":            user.Name,
				"username":        user.Username,
				"phone":           user.Phone,
				"identity_number": user.IdentityNumber,
				"password":        password,
				"email":           "invalid-email",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0).
					Return(user, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"name":            user.Name,
				"username":        user.Username,
				"phone":           user.Phone,
				"identity_number": user.IdentityNumber,
				"password":        "123",
				"email":           "invalid-email",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0).
					Return(user, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			//server := NewServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.GetUserByUsernameRow, password string) {
	username := util.RandomUsername()
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.GetUserByUsernameRow{
		//ID:       util.RandomInt(1, 1000),
		Name:     username,
		Email:    username + "@gmial.com",
		Username: username,
		Password: sql.NullString{String: hashedPassword, Valid: true},
		//Balance: sql.NullString{
		//	String: "0.00",
		//	Valid:  true,
		//},
		Phone:          "081219836581",
		IdentityNumber: "3201011411870003",
		BankCode: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
		CreatedBy: sql.NullInt64{Int64: 1, Valid: true},
		RoleID: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
		//Name_2: sql.NullString{},
	}
	return
}

func randomUserResponse(user db.GetUserByUsernameRow) dto.UserResponse {
	var roles []dto.RoleUser

	role := dto.RoleUser{
		ID:     1,
		RoleID: 1,
		UserID: user.ID,
	}

	roles = append(roles, role)

	userResponse := dto.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Username:       user.Username,
		Balance:        user.Balance,
		Phone:          user.Phone,
		IdentityNumber: user.IdentityNumber,
		BankCode:       user.BankCode.Int64,
		Role:           roles,
	}
	return userResponse
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user dto.UserResponse) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser dto.ResponseDefault
	var responseUser dto.UserResponse
	err = json.Unmarshal(data, &gotUser)

	responseUserStr, err := json.Marshal(gotUser.Data)
	require.NoError(t, err)

	err = json.Unmarshal(responseUserStr, &responseUser)

	require.NoError(t, err)
	require.Equal(t, user.ID, responseUser.ID)
	require.Equal(t, user.Name, responseUser.Name)
	require.Equal(t, user.Username, responseUser.Username)
	require.Equal(t, user.Phone, responseUser.Phone)
	require.Equal(t, user.Email, responseUser.Email)
	//require.Equal(t, user.IdentityNumber, gotUser.IdentityNumber)
	//require.Equal(t, user.Balance.String, gotUser.Balance.String)
	//require.Empty(t, gotUser.Password.String)
}

//func TestServer_getUser(t *testing.T) {
//	type fields struct {
//		store      db.Store
//		TokenMaker token.Maker
//		Router     *gin.Engine
//		config     util.Config
//	}
//	type args struct {
//		ctx *gin.Context
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			server := &Server{
//				store:      tt.fields.store,
//				TokenMaker: tt.fields.TokenMaker,
//				Router:     tt.fields.Router,
//				config:     tt.fields.config,
//			}
//			server.getUser(tt.args.ctx)
//		})
//	}
//}
