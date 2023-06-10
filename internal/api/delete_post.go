package api

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/post-service/internal/storage"
	"github.com/v1tbrah/post-service/ppbapi"
)

func (a *API) DeletePost(ctx context.Context, req *ppbapi.DeletePostRequest) (*ppbapi.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	userID, err := a.storage.DeletePost(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrPostNotFoundByID) {
			return nil, status.Error(codes.NotFound, ppbapi.ErrPostNotFoundByID.Error())
		}

		log.Error().Err(err).Int64("id", req.GetId()).Msg("storage.DeletePost")
		return nil, status.Error(codes.Internal, err.Error())
	}

	go a.msgSender.SendMsgPostDeleted(req.GetId(), userID)

	return &ppbapi.Empty{}, nil
}
