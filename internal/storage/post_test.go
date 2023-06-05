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
	ctx := context.Background()

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

			id, err := s.CreatePost(ctx, tt.post)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			postFromDB := model.Post{}
			row := s.db.QueryRow(fmt.Sprintf("SELECT user_id, created_at, description FROM table_post WHERE id=%d", id))
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
	ctx := context.Background()
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
			s := tHelperInitEmptyDB(t)

			createPostStmt, err := s.db.Prepare(`
INSERT INTO table_post (user_id, description, created_at)
VALUES ($1, $2, $3) RETURNING id
`)

			if err != nil {
				t.Fatalf("prepare create post stmt: %v", err)
			}
			res := createPostStmt.QueryRow(tt.input.UserID, tt.input.Description, tt.input.CreatedAt)
			var id int64
			if err = res.Scan(&id); err != nil {
				t.Fatalf("scan created post id: %v", err)
			}

			deletedPostUserID, err := s.DeletePost(ctx, id)
			if err != nil {
				t.Fatalf("delete post: %v", err)
			}

			if deletedPostUserID != tt.input.UserID {
				t.Fatalf("deleted post user id: %d != %d", deletedPostUserID, tt.input.UserID)
			}

			var count int64
			row := s.db.QueryRow(fmt.Sprintf("SELECT COUNT(id) FROM table_post WHERE id=%d", id))
			if err = row.Scan(&count); err != nil {
				t.Fatalf("scan count posts: %v", err)
			}
			assert.Equalf(t, int64(0), count, "count posts")
		})
	}
}

func TestStorage_GetPost(t *testing.T) {
	ctx := context.Background()

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

			createPostStmt, err := s.db.Prepare(`
INSERT INTO table_post (user_id, description, created_at)
VALUES ($1, $2, $3) RETURNING id
`)

			if err != nil {
				t.Fatalf("prepare create post stmt: %v", err)
			}
			res := createPostStmt.QueryRow(tt.input.UserID, tt.input.Description, tt.input.CreatedAt)
			var id int64
			if err = res.Scan(&id); err != nil {
				t.Fatalf("scan created post id: %v", err)
			}

			postFromDB, err := s.GetPost(ctx, id)
			tt.input.ID = postFromDB.ID
			require.NoError(t, err)
			assert.Equal(t, tt.input, postFromDB)
		})
	}
}

func TestStorage_GetPostsByUserID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		userID         int64
		input          []model.Post
		expectedOutput []model.Post
	}{
		{
			name:   "simple test",
			userID: 1,
			input: []model.Post{
				{
					UserID:      1,
					Description: "testDescription",
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
			},
			expectedOutput: []model.Post{
				{
					UserID:      1,
					Description: "testDescription",
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
			},
		},
		{
			name:   "2 on input user, one on another user",
			userID: 1,
			input: []model.Post{
				{
					UserID:      1,
					Description: "testDescription",
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
				{
					UserID:      1,
					Description: "testDescription2",
					CreatedAt:   time.Unix(20, 0).UTC(),
				},
				{
					UserID:      2,
					Description: "testDescription3",
					CreatedAt:   time.Unix(30, 0).UTC(),
				},
			},
			expectedOutput: []model.Post{
				{
					UserID:      1,
					Description: "testDescription",
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
				{
					UserID:      1,
					Description: "testDescription2",
					CreatedAt:   time.Unix(20, 0).UTC(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tHelperInitEmptyDB(t)

			for _, post := range tt.input {
				createPostStmt, err := s.db.Prepare(`
INSERT INTO table_post (user_id, description, created_at)
VALUES ($1, $2, $3) RETURNING id
`)
				if err != nil {
					t.Fatalf("prepare create post stmt: %v", err)
				}
				res := createPostStmt.QueryRow(post.UserID, post.Description, post.CreatedAt)
				var id int64
				if err = res.Scan(&id); err != nil {
					t.Fatalf("scan created post id: %v", err)
				}
			}

			postsFromDB, err := s.GetPostsByUserID(ctx, tt.userID)
			for i := range postsFromDB {
				postsFromDB[i].ID = 0 // for check equality with only input fields
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedOutput, postsFromDB)
		})
	}
}
