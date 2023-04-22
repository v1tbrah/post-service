package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.com/pet-pr-social-network/post-service/pbapi"
)

func (a *API) GetPost(ctx context.Context, req *pbapi.GetPostRequest) (*pbapi.GetPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, pbapi.ErrEmptyRequest.Error())
	}

	post, err := a.storage.GetPost(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, pbapi.ErrPostNotFoundByID.Error())
		}
		log.Error().Err(err).Int64("id", req.GetId()).Msg("storage.GetPost")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbapi.GetPostResponse{
		UserID:      post.UserID,
		Description: post.Description,
		HashtagsID:  post.HashtagsID,
		CreatedAt:   timestamppb.New(post.CreatedAt),
	}, nil
}
