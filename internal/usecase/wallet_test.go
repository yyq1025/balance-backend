package usecase_test

import (
	"context"
	"sync"
	"testing"
	"yyq1025/balance-backend/internal/entity"
	"yyq1025/balance-backend/internal/entity/mocks"
	"yyq1025/balance-backend/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func wallet(t *testing.T) (entity.WalletUseCase, *mocks.MockWalletRepository, *mocks.MockWalletEthAPI) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := mocks.NewMockWalletRepository(mockCtl)
	ethAPI := mocks.NewMockWalletEthAPI(mockCtl)
	mockWalletUseCase := usecase.NewWalletUseCase(repo, ethAPI)

	return mockWalletUseCase, repo, ethAPI
}

func TestAddOne(t *testing.T) {
	t.Parallel()

	walletUseCase, repo, ethAPI := wallet(t)

	var mu sync.Mutex

	tests := []test{
		{
			name: "Success",
			mock: func() {
				ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(1.1, nil)
				repo.EXPECT().AddOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, wallet *entity.Wallet) error {
						wallet.ID = 1
						return nil
					},
				)
			},
			res: entity.Balance{
				Wallet: entity.Wallet{
					ID:     1,
					UserID: "1",
				},
				Symbol:  "ETH",
				Balance: 1.1,
			},
			err: nil,
		},
		{
			name: "GetSymbol error",
			mock: func() {
				ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("", errInternalServErr)
			},
			res: entity.Balance{},
			err: entity.ErrAddWallet,
		},
		{
			name: "GetBalance error",
			mock: func() {
				ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(float64(0), errInternalServErr)
			},
			res: entity.Balance{},
			err: entity.ErrAddWallet,
		},
		{
			name: "Repository error",
			mock: func() {
				ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(1.1, nil)
				repo.EXPECT().AddOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{})).Return(errInternalServErr)
			},
			res: entity.Balance{},
			err: entity.ErrAddWallet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mu.Lock()
			tt.mock()
			res, err := walletUseCase.AddOne(context.Background(), &entity.Wallet{UserID: "1"})
			mu.Unlock()

			require.Equal(t, tt.res, res)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestGetOne(t *testing.T) {
	t.Parallel()

	walletUseCase, repo, ethAPI := wallet(t)

	var mu sync.Mutex

	tests := []test{
		{
			name: "Success",
			mock: func() {
				repo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, condition entity.Wallet, wallet *entity.Wallet) error {
						wallet.ID = condition.ID
						wallet.UserID = condition.UserID
						return nil
					},
				)
				ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(1.1, nil)
			},
			res: entity.Balance{
				Wallet: entity.Wallet{
					ID:     1,
					UserID: "1",
				},
				Symbol:  "ETH",
				Balance: 1.1,
			},
			err: nil,
		},
		{
			name: "Repository error",
			mock: func() {
				repo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).Return(errInternalServErr)
			},
			res: entity.Balance{},
			err: entity.ErrGetBalance,
		},
		{
			name: "GetSymbol error",
			mock: func() {
				repo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, condition entity.Wallet, wallet *entity.Wallet) error {
						wallet.ID = condition.ID
						wallet.UserID = condition.UserID
						return nil
					},
				)
				ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("", errInternalServErr)
			},
			res: entity.Balance{},
			err: entity.ErrGetBalance,
		},
		{
			name: "GetBalance error",
			mock: func() {
				repo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, condition entity.Wallet, wallet *entity.Wallet) error {
						wallet.ID = condition.ID
						wallet.UserID = condition.UserID
						return nil
					},
				)
				ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(float64(0), errInternalServErr)
			},
			res: entity.Balance{},
			err: entity.ErrGetBalance,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mu.Lock()
			tt.mock()
			res, err := walletUseCase.GetOne(context.Background(), entity.Wallet{ID: 1, UserID: "1"})
			mu.Unlock()

			require.Equal(t, tt.res, res)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestGetManyWithPagination(t *testing.T) {
	t.Parallel()

	walletUseCase, repo, ethAPI := wallet(t)

	var mu sync.Mutex

	tests := []struct {
		test
		res1 any
	}{
		{
			test: test{
				name: "First page success",
				mock: func() {
					repo.EXPECT().GetManyWithPagination(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&[]entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Pagination{})).DoAndReturn(
						func(ctx context.Context, condition entity.Wallet, wallets *[]entity.Wallet, pagination *entity.Pagination) error {
							*wallets = []entity.Wallet{
								{ID: 2, UserID: condition.UserID},
								{ID: 1, UserID: condition.UserID},
							}
							return nil
						})
					ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil).Times(2)
					ethAPI.EXPECT().GetBalance(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(1.1, nil).Times(2)
				},
				res: []entity.Balance{
					{
						Wallet: entity.Wallet{
							ID:     2,
							UserID: "1",
						},
						Symbol:  "ETH",
						Balance: 1.1,
					},
					{
						Wallet: entity.Wallet{
							ID:     1,
							UserID: "1",
						},
						Symbol:  "ETH",
						Balance: 1.1,
					},
				},
				err: nil,
			},
			res1: &entity.Pagination{
				IDLte:    2,
				Page:     1,
				PageSize: 2,
			},
		},
		{
			test: test{
				name: "Last page success",
				mock: func() {
					repo.EXPECT().GetManyWithPagination(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&[]entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Pagination{})).DoAndReturn(
						func(ctx context.Context, condition entity.Wallet, wallets *[]entity.Wallet, pagination *entity.Pagination) error {
							*wallets = []entity.Wallet{
								{ID: 1, UserID: condition.UserID},
							}
							return nil
						})
					ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil)
					ethAPI.EXPECT().GetBalance(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(1.1, nil)
				},
				res: []entity.Balance{
					{
						Wallet: entity.Wallet{
							ID:     1,
							UserID: "1",
						},
						Symbol:  "ETH",
						Balance: 1.1,
					},
				},
				err: nil,
			},
			res1: (*entity.Pagination)(nil),
		},
		{
			test: test{
				name: "Repository error",
				mock: func() {
					repo.EXPECT().GetManyWithPagination(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&[]entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Pagination{})).Return(errInternalServErr)
				},
				res: []entity.Balance(nil),
				err: entity.ErrFindWallet,
			},
			res1: (*entity.Pagination)(nil),
		},
		{
			test: test{
				name: "GetSymbol error",
				mock: func() {
					repo.EXPECT().GetManyWithPagination(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&[]entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Pagination{})).DoAndReturn(
						func(ctx context.Context, condition entity.Wallet, wallets *[]entity.Wallet, pagination *entity.Pagination) error {
							*wallets = []entity.Wallet{
								{ID: 2, UserID: condition.UserID},
								{ID: 1, UserID: condition.UserID},
							}
							return nil
						})
					ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("", errInternalServErr).Times(2)
				},
				res: []entity.Balance{
					{
						Wallet: entity.Wallet{
							ID:     2,
							UserID: "1",
						},
						Balance: -1,
					},
					{
						Wallet: entity.Wallet{
							ID:     1,
							UserID: "1",
						},
						Balance: -1,
					},
				},
				err: nil,
			},
			res1: &entity.Pagination{
				IDLte:    2,
				Page:     1,
				PageSize: 2,
			},
		},
		{
			test: test{
				name: "GetBalance error",
				mock: func() {
					repo.EXPECT().GetManyWithPagination(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{}), gomock.AssignableToTypeOf(&[]entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Pagination{})).DoAndReturn(
						func(ctx context.Context, condition entity.Wallet, wallets *[]entity.Wallet, pagination *entity.Pagination) error {
							*wallets = []entity.Wallet{
								{ID: 1, UserID: condition.UserID},
							}
							return nil
						})
					ethAPI.EXPECT().GetSymbol(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil)
					ethAPI.EXPECT().GetBalance(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(float64(0), errInternalServErr)
				},
				res: []entity.Balance{
					{
						Wallet: entity.Wallet{
							ID:     1,
							UserID: "1",
						},
						Balance: -1,
					},
				},
				err: nil,
			},
			res1: (*entity.Pagination)(nil),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mu.Lock()
			tt.mock()
			res, res1, err := walletUseCase.GetManyWithPagination(context.Background(), entity.Wallet{UserID: "1"}, &entity.Pagination{PageSize: 2})
			mu.Unlock()

			require.Equal(t, tt.res, res)
			require.Equal(t, tt.res1, res1)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestDeleteOne(t *testing.T) {
	t.Parallel()

	walletUseCase, repo, _ := wallet(t)

	var mu sync.Mutex

	tests := []test{
		{
			name: "Success",
			mock: func() {
				repo.EXPECT().DeleteOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(nil)
			},
			err: nil,
		},
		{
			name: "Repository error",
			mock: func() {
				repo.EXPECT().DeleteOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(entity.Wallet{})).Return(errInternalServErr)
			},
			err: entity.ErrDeleteWallet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mu.Lock()
			tt.mock()
			err := walletUseCase.DeleteOne(context.Background(), entity.Wallet{})
			mu.Unlock()

			require.ErrorIs(t, err, tt.err)
		})
	}
}
