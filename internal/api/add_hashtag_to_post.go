package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) AddHashtagToPost(ctx context.Context, req *ppbapi.AddHashtagToPostRequest) (*ppbapi.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ppbapi.ErrEmptyRequest.Error())
	}

	if err := a.storage.AddHashtagToPost(ctx, req.GetPostID(), req.GetHashtagID()); err != nil {
		log.Err(err).Int64("postID", req.GetPostID()).Int64("hashtagID", req.GetHashtagID()).Msg("storage.AddHashtagToPost")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ppbapi.Empty{}, nil
}
