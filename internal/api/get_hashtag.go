package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/post-service/pbapi"
)

func (a *API) GetHashtag(ctx context.Context, req *pbapi.GetHashtagRequest) (*pbapi.GetHashtagResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, pbapi.ErrEmptyRequest.Error())
	}

	hashtag, err := a.storage.GetHashtag(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, pbapi.ErrHashtagNotFoundByID.Error())
		}
		log.Err(err).Int64("id", req.GetId()).Msg("storage.GetHashtag")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbapi.GetHashtagResponse{Hashtag: &pbapi.Hashtag{Id: hashtag.ID, Name: hashtag.Name}}, nil
}
