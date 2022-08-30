package api

import (
	"fmt"
	"github.com/ariandi/ppob_go/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	createToken, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, createToken)
	request.Header.Set(AuthorizationHeaderKey, authorizationHeader)
}

//func TestAuthMiddleware(t *testing.T) {
//	testCases := []struct {
//		name          string
//		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
//		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
//	}{
//		{
//			name: "OK",
//			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
//				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, "user", time.Minute)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusOK, recorder.Code)
//			},
//		},
//		{
//			name: "NoAuthorization",
//			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusUnauthorized, recorder.Code)
//			},
//		},
//		{
//			name: "UnsupportedAuthorization",
//			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
//				addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusUnauthorized, recorder.Code)
//			},
//		},
//		{
//			name: "InvalidAuthorizationFormat",
//			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
//				addAuthorization(t, request, tokenMaker, "", "user", time.Minute)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusUnauthorized, recorder.Code)
//			},
//		},
//		{
//			name: "ExpiredToken",
//			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
//				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, "user", -time.Minute)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusUnauthorized, recorder.Code)
//			},
//		},
//	}
//
//	for i := range testCases {
//		tc := testCases[i]
//
//		t.Run(tc.name, func(t *testing.T) {
//			server := newTestServer(t, nil)
//			authPath := "/auth"
//			server.Router.GET(
//				authPath,
//				AuthMiddleware(server.TokenMaker),
//				func(ctx *gin.Context) {
//					ctx.JSON(http.StatusOK, gin.H{})
//				},
//			)
//
//			recorder := httptest.NewRecorder()
//			request, err := http.NewRequest(http.MethodGet, authPath, nil)
//			require.NoError(t, err)
//
//			tc.setupAuth(t, request, server.TokenMaker)
//			server.Router.ServeHTTP(recorder, request)
//			tc.checkResponse(t, recorder)
//		})
//	}
//}
