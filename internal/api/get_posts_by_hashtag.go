package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/v1tbrah/post-service/internal/model"
	"github.com/v1tbrah/post-service/ppbapi"
)

func (a *API) GetPostsByHashtag(ctx context.Context, req *ppbapi.GetPostsByHashtagRequest) (*ppbapi.GetPostsByHashtagResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	posts, err := a.storage.GetPostsByHashtag(ctx, req.GetHashtagID(), model.Direction(req.GetDirection()), req.GetPostOffsetID(), req.GetLimit())
	if err != nil {
		log.Err(err).Msg("storage.GetPostsByHashtag")
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

	return &ppbapi.GetPostsByHashtagResponse{Posts: res}, nil
}
