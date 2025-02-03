package build

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/kappusuton-yon-tebaru/backend/internal/rmq"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
)

type Service struct {
	rmq *rmq.BuilderRmq
}

func NewService(rmq *rmq.BuilderRmq) *Service {
	return &Service{
		rmq,
	}
}

func (s *Service) BuildImage(ctx context.Context, req BuildRequest) (string, *werror.WError) {
	req.Id = uuid.New().String()

	bs, err := json.Marshal(req)
	if err != nil {
		return "", nil
	}

	if err := s.rmq.Publish(ctx, bs); err != nil {
		return "", werror.NewFromError(err)
	}

	return req.Id, nil
}
