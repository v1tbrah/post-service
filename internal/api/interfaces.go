package api

import (
	"context"

	"gitlab.com/pet-pr-social-network/post-service/internal/model"
)

//go:generate mockery --name Storage
type Storage interface {
	CreatePost(ctx context.Context, post model.Post) (id int64, err error)
	DeletePost(ctx context.Context, id int64) error
	GetPost(ctx context.Context, id int64) (post model.Post, err error)
	GetPostsByHashtag(ctx context.Context, hashtagID int64, direction model.Direction, postOffsetID, limit int64) (posts []model.Post, err error)
	CreateHashtag(ctx context.Context, hashtag model.Hashtag) (id int64, err error)
	GetHashtag(ctx context.Context, id int64) (hashtag model.Hashtag, err error)
	AddHashtagToPost(ctx context.Context, postID, hashtagID int64) error
}

//go:generate mockery --name PostMsgSender
type PostMsgSender interface {
	SendMsgPostCreated(post model.Post)
	SendMsgPostDeleted(id int64)
}
