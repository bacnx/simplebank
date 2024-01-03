package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/bacnx/simplebank/db/mock"
	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTransfer(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()
	entry1 := randomEntry()
	entry2 := randomEntry()
	transfer := randomTransfer()

	currency := util.USD
	account1.Currency = currency
	account2.Currency = currency

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		fromAccountID int64
		toAccountID   int64
		amount        int64
		currency      string
		buildStub     func(*mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			amount:        10,
			currency:      currency,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				transferTxParams := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        10,
				}
				transferTxResult := db.TransferTxResult{
					Transfer:    transfer,
					FromAccount: account1,
					ToAccount:   account2,
					FromEntry:   entry1,
					ToEntry:     entry2,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), transferTxParams).
					Times(1).
					Return(transferTxResult, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				transferTxResult := db.TransferTxResult{
					Transfer:    transfer,
					FromAccount: account1,
					ToAccount:   account2,
					FromEntry:   entry1,
					ToEntry:     entry2,
				}
				requireBodyMatchTransferResult(t, recorder.Body, transferTxResult)
			},
		},
		{
			name:          "InvalidCurrency",
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			amount:        10,
			currency:      "InvalidCurrency",
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "TransferTxError",
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			amount:        10,
			currency:      currency,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				transferTxParams := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        10,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), transferTxParams).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "FromAccountNotFound",
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			amount:        10,
			currency:      currency,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:          "ToAccountNotFound",
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			amount:        10,
			currency:      currency,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:          "GetFromAccountError",
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			amount:        10,
			currency:      currency,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "GetToAccountError",
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			amount:        10,
			currency:      currency,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "FromAccountCurrencyMismatch",
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			amount:        10,
			currency:      util.EUR,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			body := createTransferRequest{
				FromAccountID: tc.fromAccountID,
				ToAccountID:   tc.toAccountID,
				Amount:        tc.amount,
				Currency:      tc.currency,
			}
			jsonBody, err := json.Marshal(body)
			reader := bytes.NewBuffer(jsonBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/transfer", reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomTransfer() db.Transfer {
	return db.Transfer{
		ID:            util.RandomInt(1, 1000),
		FromAccountID: util.RandomInt(1, 1000),
		ToAccountID:   util.RandomInt(1, 1000),
		Amount:        util.RandomMoney(),
	}
}

func randomEntry() db.Entry {
	return db.Entry{
		ID:        util.RandomInt(1, 1000),
		AccountID: util.RandomInt(1, 1000),
		Amount:    util.RandomMoney(),
	}
}

func requireBodyMatchTransferResult(t *testing.T, body *bytes.Buffer, result db.TransferTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResult db.TransferTxResult
	err = json.Unmarshal(data, &gotResult)
	require.NoError(t, err)
	require.Equal(t, result, gotResult)
}
