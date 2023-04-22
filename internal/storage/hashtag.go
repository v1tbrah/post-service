package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.com/pet-pr-social-network/post-service/internal/model"
)

func (s *Storage) CreateHashtag(ctx context.Context, hashtag model.Hashtag) (id int64, err error) {
	row := s.stmtHashtag.stmtCreateHashtag.QueryRowContext(ctx, hashtag.Name)
	if err = row.Scan(&id); err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == pgerrcode.UniqueViolation {
			return -1, ErrHashtagAlreadyExists
		}

		return -1, fmt.Errorf("scan created hashtag id: %w", err)
	}
	if row.Err() != nil {
		return -1, fmt.Errorf("check scan err: %w", row.Err())
	}

	return id, nil
}

func (s *Storage) GetHashtag(ctx context.Context, id int64) (hashtag model.Hashtag, err error) {
	row := s.stmtHashtag.stmtGetHashtag.QueryRowContext(ctx, id)
	if err = row.Scan(&hashtag.ID, &hashtag.Name); err != nil {
		return hashtag, fmt.Errorf("scan get hashtag by id: %w", err)
	}
	if row.Err() != nil {
		return hashtag, fmt.Errorf("check scan err: %w", row.Err())
	}

	return hashtag, nil
}
