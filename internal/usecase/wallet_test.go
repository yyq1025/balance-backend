package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"
	"yyq1025/balance-backend/internal/entity"
	"yyq1025/balance-backend/internal/entity/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()

func wallet(t *testing.T) (*entity.WalletUseCase, *mocks.MockWalletRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	rdb, _ := redismock.NewClientMock()
	mockWalletRepo := mocks.NewMockWalletRepository(mockCtl)
	mockWalletUseCase := NewWalletUseCase(mockWalletRepo, rdb)

	return &mockWalletUseCase, mockWalletRepo
}

func TestAddOne(t *testing.T) {
	t.Parallel()

	wallet, mockWalletRepo := wallet(t)

	tests := []struct {
		name     string
		ctx      context.Context
		wallet   *entity.Wallet
		mock     func()
		validate func(entity.Balance) bool
		err      error
	}{
		{
			name: "Ethereum zero address",
			ctx:  context.Background(),
			wallet: &entity.Wallet{
				Network: entity.Network{
					URL:    "https://eth.public-rpc.com",
					Symbol: "ETH",
				}},
			mock: func() {
				mockWalletRepo.EXPECT().AddOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, wallet *entity.Wallet) error {
						wallet.ID = 1
						return nil
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return balance.ID == 1 && balance.Symbol == "ETH" && balance.Balance > 0
			},
			err: nil,
		},
		{
			name: "DAI zero address",
			ctx:  context.Background(),
			wallet: &entity.Wallet{
				Token: common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
				Network: entity.Network{
					URL: "https://eth.public-rpc.com",
				}},
			mock: func() {
				mockWalletRepo.EXPECT().AddOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, wallet *entity.Wallet) error {
						wallet.ID = 2
						return nil
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return balance.ID == 2 && balance.Symbol == "DAI" && balance.Balance > 0
			},
			err: nil,
		},
		{
			name: "Repository error",
			ctx:  context.Background(),
			wallet: &entity.Wallet{
				Network: entity.Network{
					URL:    "https://eth.public-rpc.com",
					Symbol: "ETH",
				}},
			mock: func() {
				mockWalletRepo.EXPECT().AddOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, wallet *entity.Wallet) error {
						return errInternalServErr
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return true
			},
			err: entity.ErrAddWallet,
		},
		{
			name: "ETH rpc error",
			ctx:  context.Background(),
			wallet: &entity.Wallet{
				Network: entity.Network{
					URL: "abc",
				},
			},
			mock: func() {},
			validate: func(balance entity.Balance) bool {
				return true
			},
			err: entity.ErrAddWallet,
		},
		{
			name: "Token rpc error",
			ctx:  context.Background(),
			wallet: &entity.Wallet{
				Token: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
				Network: entity.Network{
					URL: "abc",
				},
			},
			mock: func() {},
			validate: func(balance entity.Balance) bool {
				return true
			},
			err: entity.ErrAddWallet,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			test.mock()
			balance, err := (*wallet).AddOne(test.ctx, test.wallet)
			require.ErrorIs(t, err, test.err)
			require.True(t, test.validate(balance))
		})
	}
}

func TestGetOne(t *testing.T) {
	t.Parallel()

	wallet, mockWalletRepo := wallet(t)

	timeout, cancel := context.WithDeadline(context.Background(), time.Now())
	defer cancel()

	tests := []struct {
		name      string
		ctx       context.Context
		condition *entity.Wallet
		mock      func()
		validate  func(entity.Balance) bool
		err       error
	}{
		{
			name: "Ethereum zero address",
			ctx:  context.Background(),
			condition: &entity.Wallet{
				ID:     1,
				UserID: "1",
			},
			mock: func() {
				mockWalletRepo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, condition, wallet *entity.Wallet) error {
						*wallet = entity.Wallet{
							ID:     1,
							UserID: "1",
							Network: entity.Network{
								URL:    "https://eth.public-rpc.com",
								Symbol: "ETH",
							},
						}
						return nil
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return balance.ID == 1 && balance.Symbol == "ETH" && balance.Balance > 0
			},
			err: nil,
		},
		{
			name: "DAI zero address",
			ctx:  context.Background(),
			condition: &entity.Wallet{
				ID:     2,
				UserID: "1",
			},
			mock: func() {
				mockWalletRepo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, condition, wallet *entity.Wallet) error {
						*wallet = entity.Wallet{
							ID:     2,
							UserID: "1",
							Token:  common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
							Network: entity.Network{
								URL: "https://eth.public-rpc.com",
							},
						}
						return nil
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return balance.ID == 2 && balance.Symbol == "DAI" && balance.Balance > 0
			},
			err: nil,
		},
		{
			name: "Repository error",
			ctx:  context.Background(),
			condition: &entity.Wallet{
				ID:     3,
				UserID: "1",
			},
			mock: func() {
				mockWalletRepo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, condition, wallet *entity.Wallet) error {
						return errInternalServErr
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return true
			},
			err: entity.ErrGetBalance,
		},
		{
			name: "RPC timeout error",
			ctx:  timeout,
			condition: &entity.Wallet{
				ID:     4,
				UserID: "1",
			},
			mock: func() {
				mockWalletRepo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, condition, wallet *entity.Wallet) error {
						*wallet = entity.Wallet{
							ID:     4,
							UserID: "1",
							Network: entity.Network{
								URL:    "https://eth.public-rpc.com",
								Symbol: "ETH",
							},
						}
						return nil
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return true
			},
			err: entity.ErrGetBalance,
		},
		{
			name: "RPC token timeout error",
			ctx:  timeout,
			condition: &entity.Wallet{
				ID:     5,
				UserID: "1",
			},
			mock: func() {
				mockWalletRepo.EXPECT().GetOne(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&entity.Wallet{}), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, condtion, wallet *entity.Wallet) error {
						*wallet = entity.Wallet{
							ID:     5,
							UserID: "1",
							Token:  common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
							Network: entity.Network{
								URL: "https://eth.public-rpc.com",
							},
						}
						return nil
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return true
			},
			err: entity.ErrGetBalance,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := (*wallet).GetOne(tt.ctx, tt.condition)
			require.ErrorIs(t, err, tt.err)
			require.True(t, tt.validate(res))
		})
	}

}
