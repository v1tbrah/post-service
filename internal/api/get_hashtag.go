package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) GetHashtag(ctx context.Context, req *ppbapi.GetHashtagRequest) (*ppbapi.GetHashtagResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	hashtag, err := a.storage.GetHashtag(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, ppbapi.ErrHashtagNotFoundByID.Error())
		}
		log.Err(err).Int64("id", req.GetId()).Msg("storage.GetHashtag")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ppbapi.GetHashtagResponse{Hashtag: &ppbapi.Hashtag{Id: hashtag.ID, Name: hashtag.Name}}, nil
}
