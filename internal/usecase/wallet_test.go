package usecase

import (
	"context"
	"testing"
	"yyq1025/balance-backend/internal/entity"
	"yyq1025/balance-backend/internal/entity/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

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
				mockWalletRepo.EXPECT().AddOne(context.Background(), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
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
				mockWalletRepo.EXPECT().AddOne(context.Background(), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
					func(ctx context.Context, wallet *entity.Wallet) error {
						wallet.ID = 1
						return nil
					},
				)
			},
			validate: func(balance entity.Balance) bool {
				return balance.ID == 1 && balance.Symbol == "DAI" && balance.Balance > 0
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
				mockWalletRepo.EXPECT().AddOne(context.Background(), gomock.AssignableToTypeOf(&entity.Wallet{})).DoAndReturn(
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
			name:   "ETH rpc error",
			ctx:    context.Background(),
			wallet: &entity.Wallet{},
			mock:   func() {},
			validate: func(balance entity.Balance) bool {
				return true
			},
			err: entity.ErrAddWallet,
		},
		{
			name: "Token rpc error",
			ctx:  context.Background(),
			wallet: &entity.Wallet{
				Token: common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
			},
			mock: func() {},
			validate: func(balance entity.Balance) bool {
				return true
			},
			err: entity.ErrAddWallet,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			// t.Parallel()

			test.mock()
			balance, err := (*wallet).AddOne(test.ctx, test.wallet)
			require.ErrorIs(t, test.err, err)
			require.True(t, test.validate(balance))
		})
	}
}
