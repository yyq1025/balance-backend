package usecase

import (
	"context"
	"errors"
	"testing"
	"yyq1025/balance-backend/internal/entity"
	"yyq1025/balance-backend/internal/entity/mocks"

	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errInternalServErr = errors.New("cannot get networks")

func network(t *testing.T) (*entity.NetworkUseCase, *mocks.MockNetworkRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	rdb, _ := redismock.NewClientMock()
	mockNetworkRepo := mocks.NewMockNetworkRepository(mockCtl)
	mockNetworkUseCase := NewNetworkUseCase(mockNetworkRepo, rdb)

	return &mockNetworkUseCase, mockNetworkRepo
}

func TestGetAll(t *testing.T) {
	t.Parallel()

	network, mockNetworkRepo := network(t)

	tests := []struct {
		name string
		mock func()
		res  any
		err  error
	}{
		{
			name: "error",
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
		{
			name: "not cached",
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
			name: "cached",
			mock: func() {},
			res:  []entity.Network{{Name: "test"}},
			err:  nil,
		},
	}

	for _, tt := range tests {
		// tt := tt
		t.Run(tt.name, func(t *testing.T) {

			tt.mock()
			res, err := (*network).GetAll(context.Background())
			require.Equal(t, tt.res, res)
			require.ErrorIs(t, tt.err, err)
		})
	}
}
