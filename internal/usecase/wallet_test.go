package usecase_test

import (
	"context"
	"testing"

	"github.com/yyq1025/balance-backend/internal/entity"
	"github.com/yyq1025/balance-backend/internal/entity/mocks"
	"github.com/yyq1025/balance-backend/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func wallet(t *testing.T) (entity.WalletUseCase, *mocks.MockWalletRepository, *mocks.MockWalletEthAPI) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := mocks.NewMockWalletRepository(mockCtl)
	ethAPI := mocks.NewMockWalletEthAPI(mockCtl)
	walletUseCase := usecase.NewWalletUseCase(repo, ethAPI)

	return walletUseCase, repo, ethAPI
}

func TestAddOne(t *testing.T) {
	t.Parallel()

	walletUseCase, repo, ethAPI := wallet(t)

	tests := []test{
		{
			name: "Success",
			mock: func() {
				ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{UserID: "1"}).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(context.Background(), entity.Wallet{UserID: "1"}).Return(1.1, nil)
				repo.EXPECT().AddOne(context.Background(), &entity.Wallet{UserID: "1"}).DoAndReturn(
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
				ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{UserID: "1"}).Return("", errInternalServerErr)
			},
			res: entity.Balance{},
			err: entity.ErrAddWallet,
		},
		{
			name: "GetBalance error",
			mock: func() {
				ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{UserID: "1"}).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(context.Background(), entity.Wallet{UserID: "1"}).Return(float64(0), errInternalServerErr)
			},
			res: entity.Balance{},
			err: entity.ErrAddWallet,
		},
		{
			name: "Repository error",
			mock: func() {
				ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{UserID: "1"}).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(context.Background(), entity.Wallet{UserID: "1"}).Return(1.1, nil)
				repo.EXPECT().AddOne(context.Background(), &entity.Wallet{UserID: "1"}).Return(errInternalServerErr)
			},
			res: entity.Balance{},
			err: entity.ErrAddWallet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			res, err := walletUseCase.AddOne(context.Background(), &entity.Wallet{UserID: "1"})

			require.Equal(t, tt.res, res)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestGetOne(t *testing.T) {
	t.Parallel()

	walletUseCase, repo, ethAPI := wallet(t)

	tests := []test{
		{
			name: "Success",
			mock: func() {
				repo.EXPECT().GetOne(context.Background(), "1", 1, gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, userID string, id int, wallet *entity.Wallet) error {
						wallet.ID = id
						wallet.UserID = userID
						return nil
					},
				)
				ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return(1.1, nil)
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
				repo.EXPECT().GetOne(context.Background(), "1", 1, gomock.AssignableToTypeOf(&entity.Wallet{})).Return(errInternalServerErr)
			},
			res: entity.Balance{},
			err: entity.ErrGetBalance,
		},
		{
			name: "GetSymbol error",
			mock: func() {
				repo.EXPECT().GetOne(context.Background(), "1", 1, gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, userID string, id int, wallet *entity.Wallet) error {
						wallet.ID = id
						wallet.UserID = userID
						return nil
					},
				)
				ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return("", errInternalServerErr)
			},
			res: entity.Balance{},
			err: entity.ErrGetBalance,
		},
		{
			name: "GetBalance error",
			mock: func() {
				repo.EXPECT().GetOne(context.Background(), "1", 1, gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, userID string, id int, wallet *entity.Wallet) error {
						wallet.ID = id
						wallet.UserID = userID
						return nil
					},
				)
				ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return("ETH", nil)
				ethAPI.EXPECT().GetBalance(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return(float64(0), errInternalServerErr)
			},
			res: entity.Balance{},
			err: entity.ErrGetBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			res, err := walletUseCase.GetOne(context.Background(), "1", 1)

			require.Equal(t, tt.res, res)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestGetManyWithPagination(t *testing.T) {
	t.Parallel()

	walletUseCase, repo, ethAPI := wallet(t)

	tests := []struct {
		test
		res1 any
	}{
		{
			test: test{
				name: "First page success",
				mock: func() {
					repo.EXPECT().GetManyWithPagination(context.Background(), "1", &entity.Pagination{PageSize: 2}, gomock.AssignableToTypeOf(&[]entity.Wallet{})).DoAndReturn(
						func(ctx context.Context, userID string, pagination *entity.Pagination, wallets *[]entity.Wallet) error {
							*wallets = []entity.Wallet{
								{ID: 2, UserID: userID},
								{ID: 1, UserID: userID},
							}
							return nil
						})
					ethAPI.EXPECT().GetSymbol(context.Background(), gomock.AssignableToTypeOf(entity.Wallet{})).Return("ETH", nil).Times(2)
					ethAPI.EXPECT().GetBalance(context.Background(), gomock.AssignableToTypeOf(entity.Wallet{})).Return(1.1, nil).Times(2)
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
					repo.EXPECT().GetManyWithPagination(context.Background(), "1", &entity.Pagination{PageSize: 2}, gomock.AssignableToTypeOf(&[]entity.Wallet{})).DoAndReturn(
						func(ctx context.Context, userID string, pagination *entity.Pagination, wallets *[]entity.Wallet) error {
							*wallets = []entity.Wallet{
								{ID: 1, UserID: userID},
							}
							return nil
						})
					ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return("ETH", nil)
					ethAPI.EXPECT().GetBalance(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return(1.1, nil)
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
					repo.EXPECT().GetManyWithPagination(context.Background(), "1", &entity.Pagination{PageSize: 2}, gomock.AssignableToTypeOf(&[]entity.Wallet{})).Return(errInternalServerErr)
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
					repo.EXPECT().GetManyWithPagination(context.Background(), "1", &entity.Pagination{PageSize: 2}, gomock.AssignableToTypeOf(&[]entity.Wallet{})).DoAndReturn(
						func(ctx context.Context, userID string, pagination *entity.Pagination, wallets *[]entity.Wallet) error {
							*wallets = []entity.Wallet{
								{ID: 2, UserID: userID},
								{ID: 1, UserID: userID},
							}
							return nil
						})
					ethAPI.EXPECT().GetSymbol(context.Background(), gomock.AssignableToTypeOf(entity.Wallet{})).Return("", errInternalServerErr).Times(2)
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
					repo.EXPECT().GetManyWithPagination(context.Background(), "1", &entity.Pagination{PageSize: 2}, gomock.AssignableToTypeOf(&[]entity.Wallet{})).DoAndReturn(
						func(ctx context.Context, userID string, pagination *entity.Pagination, wallets *[]entity.Wallet) error {
							*wallets = []entity.Wallet{
								{ID: 1, UserID: userID},
							}
							return nil
						})
					ethAPI.EXPECT().GetSymbol(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return("ETH", nil)
					ethAPI.EXPECT().GetBalance(context.Background(), entity.Wallet{ID: 1, UserID: "1"}).Return(float64(0), errInternalServerErr)
				},
				res: []entity.Balance{
					{
						Wallet: entity.Wallet{
							ID:     1,
							UserID: "1",
						},
						Symbol:  "ETH",
						Balance: -1,
					},
				},
				err: nil,
			},
			res1: (*entity.Pagination)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			res, res1, err := walletUseCase.GetManyWithPagination(context.Background(), "1", &entity.Pagination{PageSize: 2})

			require.Equal(t, tt.res, res)
			require.Equal(t, tt.res1, res1)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestDeleteOne(t *testing.T) {
	t.Parallel()

	walletUseCase, repo, _ := wallet(t)

	tests := []test{
		{
			name: "Success",
			mock: func() {
				repo.EXPECT().DeleteOne(context.Background(), "", 0).Return(nil)
			},
			err: nil,
		},
		{
			name: "Repository error",
			mock: func() {
				repo.EXPECT().DeleteOne(context.Background(), "", 0).Return(errInternalServerErr)
			},
			err: entity.ErrDeleteWallet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := walletUseCase.DeleteOne(context.Background(), "", 0)

			require.ErrorIs(t, err, tt.err)
		})
	}
}
