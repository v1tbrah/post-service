package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/post-service/pbapi"
)

func (a *API) AddHashtagToPost(ctx context.Context, req *pbapi.AddHashtagToPostRequest) (*pbapi.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, pbapi.ErrEmptyRequest.Error())
	}

	if err := a.storage.AddHashtagToPost(ctx, req.GetPostID(), req.GetHashtagID()); err != nil {
		log.Err(err).Int64("postID", req.GetPostID()).Int64("hashtagID", req.GetHashtagID()).Msg("storage.AddHashtagToPost")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbapi.Empty{}, nil
}
