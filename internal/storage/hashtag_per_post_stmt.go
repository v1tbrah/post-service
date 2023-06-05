package storage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type hashtagPerPost struct {
	addHashtagToPost *sql.Stmt
}

func (shpp *hashtagPerPost) prepare(db *sql.DB) (err error) {
	const addHashtagToPost = `
		INSERT INTO table_hashtag_per_post (post_id, hashtag_id)
		VALUES ($1, $2)
`

	if shpp.addHashtagToPost, err = db.Prepare(addHashtagToPost); err != nil {
		return errors.Wrap(err, "prepare 'add hashtag to post' stmt")
	}

	return nil
}

func (shpp *hashtagPerPost) close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = shpp.addHashtagToPost.Close(); err != nil {
			err = errors.Wrap(err, "close stmt 'add hashtag to post'")
			closeEnded <- struct{}{}
			return
		}

		closeEnded <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-closeEnded:
		return err
	}
}
