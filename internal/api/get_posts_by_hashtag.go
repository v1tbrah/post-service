package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/post-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.com/pet-pr-social-network/post-service/pbapi"
)

func (a *API) GetPostsByHashtag(ctx context.Context, req *pbapi.GetPostsByHashtagRequest) (*pbapi.GetPostsByHashtagResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, pbapi.ErrEmptyRequest.Error())
	}

	posts, err := a.storage.GetPostsByHashtag(ctx, req.GetHashtagID(), model.Direction(req.GetDirection()), req.GetPostOffsetID(), req.GetLimit())
	if err != nil {
		log.Err(err).Msg("storage.GetPostsByHashtag")
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := make([]*pbapi.Post, 0, len(posts))
	for _, post := range posts {
		res = append(res, &pbapi.Post{
			Id:          post.ID,
			UserID:      post.UserID,
			Description: post.Description,
			HashtagsID:  post.HashtagsID,
			CreatedAt:   timestamppb.New(post.CreatedAt),
		})
	}

	return &pbapi.GetPostsByHashtagResponse{Posts: res}, nil
}
