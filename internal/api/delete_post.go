package api

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/post-service/internal/storage"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) DeletePost(ctx context.Context, req *ppbapi.DeletePostRequest) (*ppbapi.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	err := a.storage.DeletePost(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrPostNotFoundByID) {
			return nil, status.Error(codes.NotFound, ppbapi.ErrPostNotFoundByID.Error())
		}

		log.Error().Err(err).Int64("id", req.GetId()).Msg("storage.DeletePost")
		return nil, status.Error(codes.Internal, err.Error())
	}

	go a.msgSender.SendMsgPostDeleted(req.GetId())

	return &ppbapi.Empty{}, nil
}
