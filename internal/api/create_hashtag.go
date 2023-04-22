package api

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/post-service/internal/model"
	"gitlab.com/pet-pr-social-network/post-service/internal/storage"
	"gitlab.com/pet-pr-social-network/post-service/pbapi"
)

func (a *API) CreateHashtag(ctx context.Context, req *pbapi.CreateHashtagRequest) (*pbapi.CreateHashtagResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, pbapi.ErrEmptyRequest.Error())
	}

	reqName := strings.TrimSpace(req.GetName())
	if reqName == "" {
		return nil, status.Error(codes.InvalidArgument, pbapi.ErrEmptyName.Error())
	}

	id, err := a.storage.CreateHashtag(ctx, model.Hashtag{Name: reqName})
	if err != nil {
		if errors.Is(err, storage.ErrHashtagAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		log.Err(err).Str("name", reqName).Msg("storage.CreateHashtag")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbapi.CreateHashtagResponse{Id: id}, nil
}