//go:build with_db

package storage

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/pet-pr-social-network/post-service/internal/model"
)

func TestStorage_AddHashtagToPost(t *testing.T) {
	ctx := context.Background()

	s := tHelperInitEmptyDB(t)

	tests := []struct {
		name    string
		post    model.Post
		hashtag model.Hashtag
	}{
		{
			name: "simple test",
			post: model.Post{
				UserID:      1,
				Description: "testDescription",
				CreatedAt:   time.Unix(10, 0).UTC(),
			},
			hashtag: model.Hashtag{
				Name: "testName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s = tHelperInitEmptyDB(t)

			postID, err := s.CreatePost(ctx, tt.post)
			require.NoError(t, err)

			hashtagID, err := s.CreateHashtag(ctx, tt.hashtag)
			require.NoError(t, err)

			err = s.AddHashtagToPost(ctx, postID, hashtagID)
			require.NoError(t, err)

			row := s.db.QueryRow(fmt.Sprintf("SELECT post_id, hashtag_id FROM table_hashtag_per_post WHERE post_id=%d AND hashtag_id=%d", postID, hashtagID))
			var postIDFromDB, hashtagIDFromDB int64
			if err = row.Scan(&postIDFromDB, &hashtagIDFromDB); err != nil {
				t.Fatalf("scan: %v", err)
			}

			assert.Equal(t, postID, postIDFromDB)
			assert.Equal(t, hashtagID, hashtagIDFromDB)
		})
	}
}

func TestStorage_GetPostsByHashtag(t *testing.T) {
	// Scenario:
	// set 100 pairs post+hashtag with hashtag 'forSearch', and 100 with hashtag 'notForSearch'
	// get first 20 by hashtag 'forSearch'
	// get next 20 by hashtag 'forSearch'
	// get next 60 by hashtag 'forSearch'
	// get prev 20 by hashtag 'forSearch'

	ctx := context.Background()

	s := tHelperInitEmptyDB(t)

	var tmpForSearch, tmpNotForSearch model.Post
	hashForSearch := model.Hashtag{
		Name: "forSearch",
	}
	hashNotForSearch := model.Hashtag{
		Name: "notForSearch",
	}

	hashForSearchID, err := s.CreateHashtag(ctx, hashForSearch)
	require.NoError(t, err)
	hashNotForSearchID, err := s.CreateHashtag(ctx, hashNotForSearch)
	require.NoError(t, err)

	var tmpForSearchID, tmpNotForSearchID int64
	for i := 0; i < 100; i++ {
		tmpForSearch = model.Post{
			UserID:      1,
			Description: strconv.Itoa(i + 1),
		}
		tmpForSearchID, err = s.CreatePost(ctx, tmpForSearch)
		require.NoError(t, err)

		tmpNotForSearch = model.Post{
			UserID:      1,
			Description: strconv.Itoa(-(i + 1)),
		}
		tmpNotForSearchID, err = s.CreatePost(ctx, tmpNotForSearch)
		require.NoError(t, err)

		err = s.AddHashtagToPost(ctx, tmpForSearchID, hashForSearchID)
		require.NoError(t, err)
		err = s.AddHashtagToPost(ctx, tmpNotForSearchID, hashNotForSearchID)
		require.NoError(t, err)
	}

	first20, err := s.GetPostsByHashtag(ctx, hashForSearchID, model.First, 0, 20)
	require.NoError(t, err)
	for i := 0; i < len(first20); i++ {
		require.Equal(t, strconv.Itoa(i+1), first20[i].Description)
	}

	currLast := first20[len(first20)-1]
	next20, err := s.GetPostsByHashtag(ctx, hashForSearchID, model.Next, currLast.ID, 20)
	require.NoError(t, err)
	for i := 0; i < len(next20); i++ {
		require.Equal(t, strconv.Itoa(i+1+len(first20)), next20[i].Description)
	}

	currLast = next20[len(next20)-1]
	next60, err := s.GetPostsByHashtag(ctx, hashForSearchID, model.Next, currLast.ID, 60)
	require.NoError(t, err)
	for i := 0; i < len(next60); i++ {
		require.Equal(t, strconv.Itoa(i+1+len(first20)+len(next20)), next60[i].Description)
	}

	currLast = next60[len(next60)-1]
	prev20, err := s.GetPostsByHashtag(ctx, hashForSearchID, model.Prev, currLast.ID+1, 20)
	require.NoError(t, err)
	for i := 0; i < len(prev20); i++ {
		require.Equal(t, strconv.Itoa(len(first20)+len(next20)+len(next60)-i), prev20[i].Description)
	}
}
