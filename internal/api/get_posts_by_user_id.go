package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *API) GetPostsByUserID(ctx context.Context, req *ppbapi.GetPostsByUserIDRequest) (*ppbapi.GetPostsByUserIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	posts, err := a.storage.GetPostsByUserID(ctx, req.GetUserID())
	if err != nil {
		log.Err(err).Msg("storage.GetPostsByUserID")
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := make([]*ppbapi.Post, 0, len(posts))
	for _, post := range posts {
		res = append(res, &ppbapi.Post{
			Id:          post.ID,
			UserID:      post.UserID,
			Description: post.Description,
			HashtagsID:  post.HashtagsID,
			CreatedAt:   timestamppb.New(post.CreatedAt),
		})
	}

	return &ppbapi.GetPostsByUserIDResponse{Posts: res}, nil
}
