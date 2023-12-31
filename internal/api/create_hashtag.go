package api

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/post-service/internal/model"
	"github.com/v1tbrah/post-service/internal/storage"
	"github.com/v1tbrah/post-service/ppbapi"
)

func (a *API) CreateHashtag(ctx context.Context, req *ppbapi.CreateHashtagRequest) (*ppbapi.CreateHashtagResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	reqName := strings.TrimSpace(req.GetName())
	if reqName == "" {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyName.Error())
	}

	id, err := a.storage.CreateHashtag(ctx, model.Hashtag{Name: reqName})
	if err != nil {
		if errors.Is(err, storage.ErrHashtagAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		log.Err(err).Str("name", reqName).Msg("storage.CreateHashtag")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ppbapi.CreateHashtagResponse{Id: id}, nil
}
