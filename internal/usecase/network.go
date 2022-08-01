package usecase

import (
	"context"
	"log"

	"github.com/yyq1025/balance-backend/internal/entity"
)

type networkUseCase struct {
	repo entity.NetworkRepository
}

func NewNetworkUseCase(n entity.NetworkRepository) entity.NetworkUseCase {
	return &networkUseCase{n}
}

func (n *networkUseCase) GetAll(ctx context.Context) ([]entity.Network, error) {
	var networks []entity.Network
	if err := n.repo.GetAll(ctx, &networks); err != nil {
		log.Print(err)
		return nil, entity.ErrGetNetwork
	}
	return networks, nil
}
