package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/post-service/internal/model"
)

func (s *Storage) CreatePost(ctx context.Context, post model.Post) (id int64, err error) {
	tx, err := s.dbConn.Begin()
	if err != nil {
		return -1, err
	}

	defer func() {
		if errRollback := tx.Rollback(); errRollback != nil && errRollback != sql.ErrTxDone {
			log.Error().Err(errRollback).Msg("storage.CreatePost tx.Rollback")
		}
	}()

	row := tx.Stmt(s.stmtPost.stmtCreatePost).QueryRowContext(ctx, post.UserID, post.Description, post.CreatedAt)
	if err = row.Scan(&id); err != nil {
		return -1, fmt.Errorf("scan created post id: %w", err)
	}
	if row.Err() != nil {
		return -1, fmt.Errorf("check scan err: %w", row.Err())
	}

	for _, hid := range post.HashtagsID {
		if _, err = tx.Stmt(s.stmtHashtagPerPost.stmtAddHashtagToPost).Exec(id, hid); err != nil {
			return -1, fmt.Errorf("add hashtag (%d) to post: %w", hid, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return id, nil
}

func (s *Storage) DeletePost(ctx context.Context, id int64) (userID int64, err error) {
	if err = s.stmtPost.stmtDeletePost.QueryRowContext(ctx, id).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrPostNotFoundByID
		}
		return -1, fmt.Errorf("scan deleted post userID: %w", err)
	}

	return userID, nil
}

func (s *Storage) GetPost(ctx context.Context, id int64) (post model.Post, err error) {
	row := s.stmtPost.stmtGetPost.QueryRowContext(ctx, id)
	if err = row.Scan(&post.ID, &post.UserID, &post.Description, &post.CreatedAt); err != nil {
		return post, fmt.Errorf("scan get post by id: %w", err)
	}
	if row.Err() != nil {
		return post, fmt.Errorf("check scan err: %w", row.Err())
	}

	return post, nil
}

func (s *Storage) GetPostsByUserID(ctx context.Context, userID int64) (posts []model.Post, err error) {
	rows, err := s.stmtPost.stmtGetPostsByUserID.QueryContext(ctx, userID)
	for rows.Next() {
		var post model.Post
		if err = rows.Scan(&post.ID, &post.UserID, &post.Description, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Close(); err != nil {
		return nil, fmt.Errorf("close rows: %w", err)
	}

	return posts, nil
}
