package usecase_test

import (
	"context"
	"sync"
	"testing"

	"github.com/yyq1025/balance-backend/internal/entity"
	"github.com/yyq1025/balance-backend/internal/entity/mocks"
	"github.com/yyq1025/balance-backend/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func network(t *testing.T) (entity.NetworkUseCase, *mocks.MockNetworkRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockNetworkRepo := mocks.NewMockNetworkRepository(mockCtl)
	mockNetworkUseCase := usecase.NewNetworkUseCase(mockNetworkRepo)

	return mockNetworkUseCase, mockNetworkRepo
}

func TestGetAll(t *testing.T) {
	t.Parallel()

	network, mockNetworkRepo := network(t)

	var mu sync.Mutex

	tests := []test{
		{
			name: "Success",
			mock: func() {
				mockNetworkRepo.EXPECT().GetAll(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&[]entity.Network{})).DoAndReturn(
					func(ctx context.Context, networks *[]entity.Network) error {
						*networks = []entity.Network{{Name: "test"}}
						return nil
					},
				)
			},
			res: []entity.Network{{Name: "test"}},
			err: nil,
		},
		{
			name: "Error",
			mock: func() {
				mockNetworkRepo.EXPECT().GetAll(gomock.AssignableToTypeOf(ctxType), gomock.AssignableToTypeOf(&[]entity.Network{})).DoAndReturn(
					func(ctx context.Context, networks *[]entity.Network) error {
						return errInternalServErr
					},
				)
			},
			res: []entity.Network(nil),
			err: entity.ErrGetNetwork,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mu.Lock()
			tt.mock()
			res, err := network.GetAll(context.Background())
			mu.Unlock()

			require.Equal(t, tt.res, res)
			require.ErrorIs(t, err, tt.err)
		})
	}
}