package api

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/post-service/internal/model"
	"github.com/v1tbrah/post-service/ppbapi"
)

func (a *API) CreatePost(ctx context.Context, req *ppbapi.CreatePostRequest) (*ppbapi.CreatePostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	post := model.Post{UserID: req.GetUserID(), Description: req.GetDescription(), HashtagsID: req.GetHashtagsID(), CreatedAt: time.Now()}
	id, err := a.storage.CreatePost(ctx, post)
	if err != nil {
		log.Error().Err(err).Str("user", fmt.Sprintf("%+v", req)).Msg("storage.CreatePost")
		return nil, status.Error(codes.Internal, err.Error())
	}

	post.ID = id
	go a.msgSender.SendMsgPostCreated(post)

	return &ppbapi.CreatePostResponse{Id: id}, nil
}
