package storage

import (
	"context"
	"fmt"

	"github.com/v1tbrah/post-service/internal/model"
)

func (s *Storage) GetPostsByHashtag(ctx context.Context, hashtagID int64, direction model.Direction, postOffsetID, limit int64) (posts []model.Post, err error) {
	getPostsQuery := fmt.Sprintf(`SELECT
		hpp.post_id,
		p.user_id,
		p.description,
    	p.created_at
		FROM table_hashtag_per_post AS hpp
			INNER JOIN table_post AS p
				ON hpp.post_id = p.id
		WHERE hpp.hashtag_id=%d`, hashtagID)

	switch direction {
	case model.First:
		getPostsQuery += ` ORDER BY post_id ASC`
	case model.Next:
		getPostsQuery += fmt.Sprintf(" AND post_id > %d ORDER BY post_id ASC", postOffsetID)
	case model.Prev:
		getPostsQuery += fmt.Sprintf(" AND post_id < %d ORDER BY post_id DESC", postOffsetID)
	default:
		return nil, fmt.Errorf("get posts by hashtag: %w", model.ErrInvalidDirection)
	}
	if limit > 0 {
		getPostsQuery += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := s.db.QueryContext(ctx, getPostsQuery)
	if err != nil {
		return nil, fmt.Errorf("get posts by hashtag| hashtagID: %d | postOffsetID: %d | limit: %d: %w", hashtagID, postOffsetID, limit, err)
	}

	var post model.Post
	for rows.Next() {
		if err = rows.Scan(&post.ID, &post.UserID, &post.Description, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Storage) AddHashtagToPost(ctx context.Context, postID, hashtagID int64) error {
	if _, err := s.hashtagPerPost.addHashtagToPost.ExecContext(ctx, postID, hashtagID); err != nil {
		return fmt.Errorf("add hashatag to post| postID: %d | hashtagID: %d: %w", postID, hashtagID, err)
	}

	return nil
}
