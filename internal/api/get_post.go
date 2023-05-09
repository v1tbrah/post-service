package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *API) GetPost(ctx context.Context, req *ppbapi.GetPostRequest) (*ppbapi.GetPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	post, err := a.storage.GetPost(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, ppbapi.ErrPostNotFoundByID.Error())
		}
		log.Error().Err(err).Int64("id", req.GetId()).Msg("storage.GetPost")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ppbapi.GetPostResponse{
		UserID:      post.UserID,
		Description: post.Description,
		HashtagsID:  post.HashtagsID,
		CreatedAt:   timestamppb.New(post.CreatedAt),
	}, nil
}
