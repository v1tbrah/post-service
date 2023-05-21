//go:build with_db

package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/pet-pr-social-network/post-service/internal/model"
)

func TestStorage_CreatePost(t *testing.T) {
	s := tHelperInitEmptyDB(t)

	tests := []struct {
		name    string
		post    model.Post
		wantErr bool
	}{
		{
			name: "simple test",
			post: model.Post{
				UserID:      1,
				Description: "testDescription",
				CreatedAt:   time.Unix(10, 0).UTC(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s = tHelperInitEmptyDB(t)

			id, err := s.CreatePost(context.Background(), tt.post)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			postFromDB := model.Post{}
			row := s.dbConn.QueryRow(fmt.Sprintf("SELECT user_id, created_at, description FROM %s WHERE id=%d", s.cfg.PostTableName, id))
			if err = row.Scan(&postFromDB.UserID, &postFromDB.CreatedAt, &postFromDB.Description); err != nil {
				t.Fatalf("scan new post: %v", err)
			}
			if row.Err() != nil {
				t.Fatalf("check scan err: %v", err)
			}

			assert.Equal(t, tt.post, postFromDB)
		})
	}
}

func TestStorage_DeletePost(t *testing.T) {
	s := tHelperInitEmptyDB(t)

	tests := []struct {
		name  string
		input model.Post
	}{
		{
			name: "simple test",
			input: model.Post{
				UserID:      1,
				Description: "testDescription",
				CreatedAt:   time.Unix(10, 0).UTC(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s = tHelperInitEmptyDB(t)

			createPostStmt, err := s.dbConn.Prepare(fmt.Sprintf(`
INSERT INTO %s (user_id, description, created_at)
VALUES ($1, $2, $3) RETURNING id
`, s.cfg.PostTableName))

			if err != nil {
				t.Fatalf("prepare create post stmt: %v", err)
			}
			res := createPostStmt.QueryRow(tt.input.UserID, tt.input.Description, tt.input.CreatedAt)
			var id int64
			if err = res.Scan(&id); err != nil {
				t.Fatalf("scan created post id: %v", err)
			}

			if err = s.DeletePost(context.Background(), id); err != nil {
				t.Fatalf("delete post: %v", err)
			}

			var count int64
			row := s.dbConn.QueryRow(fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE id=%d", s.cfg.PostTableName, id))
			if err = row.Scan(&count); err != nil {
				t.Fatalf("scan count posts: %v", err)
			}
			assert.Equalf(t, int64(0), count, "count posts")
		})
	}
}

func TestStorage_GetPost(t *testing.T) {
	s := tHelperInitEmptyDB(t)

	tests := []struct {
		name  string
		input model.Post
	}{
		{
			name: "simple test",
			input: model.Post{
				UserID:      1,
				Description: "testDescription",
				CreatedAt:   time.Unix(10, 0).UTC(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s = tHelperInitEmptyDB(t)

			createPostStmt, err := s.dbConn.Prepare(fmt.Sprintf(`
INSERT INTO %s (user_id, description, created_at)
VALUES ($1, $2, $3) RETURNING id
`, s.cfg.PostTableName))

			if err != nil {
				t.Fatalf("prepare create post stmt: %v", err)
			}
			res := createPostStmt.QueryRow(tt.input.UserID, tt.input.Description, tt.input.CreatedAt)
			var id int64
			if err = res.Scan(&id); err != nil {
				t.Fatalf("scan created post id: %v", err)
			}

			postFromDB, err := s.GetPost(context.Background(), id)
			tt.input.ID = postFromDB.ID
			require.NoError(t, err)
			assert.Equal(t, tt.input, postFromDB)
		})
	}
}
