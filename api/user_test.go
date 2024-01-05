package api

//
//import (
//	"context"
//	"database/sql"
//	db "github.com/kanishkmittal55/simplebank/db/sqlc"
//	"github.com/lib/pq"
//	"golang.org/x/crypto/bcrypt"
//	"io/ioutil"
//	"testing"
//	"time"
//)
//
//import (
//	"bytes"
//	"encoding/json"
//	"github.com/golang/mock/gomock"
//	mockdb "github.com/kanishkmittal55/simplebank/db/mock"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//)
//
//
//
//func TestCreateUserAPI(t *testing.T) {
//
//	testCases := []struct {
//		name             string
//		expectedResponse UserResponse
//		buildStubs       func(store *mockdb.MockStore)
//		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder)
//		requestBody      createUserRequest
//	}{
//		{
//			name: "OK",
//			expectedResponse: UserResponse{
//				Username:         "user1",
//				FullName:         "User One",
//				Email:            "user1@example.com",
//				PasswordChangeAt: time.Now(),
//				CreatedAt:        time.Now(),
//			},
//
//			requestBody: createUserRequest{
//				Username: "user1",
//				Password: "password",
//				FullName: "User One",
//				Email:    "user1@example.com",
//			},
//
//			buildStubs: func(store *mockdb.MockStore) {
//
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					DoAndReturn(func(_ context.Context, params db.CreateUserParams) (UserResponse, error) {
//						// Inside DoAndReturn, you can make assertions or run custom checks.
//
//						// Check if Username, FullName, and Email match the expected values.
//						require.Equal(t, "user1", params.Username)
//						require.Equal(t, "User One", params.FullName)
//						require.Equal(t, "user1@example.com", params.Email)
//
//						// Now, check that the hashed password is valid using bcrypt.
//						err := bcrypt.CompareHashAndPassword([]byte(params.HashedPassword), []byte("password"))
//						require.NoError(t, err)
//
//						// If the password matches, return the expected user object and nil error.
//						return User, nil
//					}).
//					Times(1).
//					Return(User, nil)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusOK, recorder.Code)
//				requireBodyMatchUser(t, recorder.Body, User)
//			},
//		},
//		{
//			name: "DuplicateUsername",
//			user: User,
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, &pq.Error{Code: "23505"})
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusForbidden, recorder.Code)
//			},
//		},
//		{
//			name: "BadRequest",
//			user: User,
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//			},
//		},
//		{
//			name: "InternalError",
//			user: User,
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, sql.ErrConnDone)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusInternalServerError, recorder.Code)
//			},
//		},
//	}
//
//	for i := range testCases {
//		tc := testCases[i]
//
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			store := mockdb.NewMockStore(ctrl)
//			tc.buildStubs(store)
//
//			// start test server and send request
//			server := newTestServer(t, store)
//			recorder := httptest.NewRecorder()
//
//			// marshal body data to JSON
//			data, err := json.Marshal(tc.requestBody)
//			require.NoError(t, err)
//
//			// create a new HTTP request with the user data
//			url := "/users"
//			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
//			require.NoError(t, err)
//
//			// record the response
//			server.router.ServeHTTP(recorder, req)
//			tc.checkResponse(t, recorder)
//
//		})
//	}
//}
//
////func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
////	// Firstly we call ioutil.ReadAll() to read all data from the response Body
////	data, err := ioutil.ReadAll(body)
////	require.NoError(t, err)
////
////	// Then we declare a new GotAccount Variable to store the account object we got from the response body data.
////	var gotUser db.User
////	// Now let's call json.Unmarshal to UNMARSHAL THE DATA to the got Account Object
////	err = json.Unmarshal(data, &gotUser)
////	require.NoError(t, err)
////	require.Equal(t, user, gotUser)
////
////}
//
//func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
//	data, err := ioutil.ReadAll(body)
//	require.NoError(t, err)
//
//	var gotUser UserResponse
//	err = json.Unmarshal(data, &gotUser)
//	require.NoError(t, err)
//
//	//expectedUser := newUserResponse(user, accessToken)
//	require.Equal(t, user, gotUser)
//}
