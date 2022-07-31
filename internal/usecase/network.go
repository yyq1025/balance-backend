package usecase

import (
	"context"
	"time"
	"yyq1025/balance-backend/internal/entity"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type networkUseCase struct {
	networkRepo entity.NetworkRepository
	cache       *cache.Cache
}

func NewNetworkUseCase(n entity.NetworkRepository, rdb *redis.Client) entity.NetworkUseCase {
	return &networkUseCase{
		networkRepo: n,
		cache: cache.New(&cache.Options{
			Redis:      rdb,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})}
}

func (n *networkUseCase) GetAll(ctx context.Context) (networks []entity.Network, err error) {
	if err = n.cache.Get(ctx, "networks", &networks); err != nil {
		if err = n.networkRepo.GetAll(ctx, &networks); err != nil {
			err = entity.ErrGetNetwork
			return
		}
		_ = n.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   "networks",
			Value: networks,
			TTL:   time.Hour,
			SetNX: true,
		})
	}

	return networks, nil
}
