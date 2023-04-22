package storage

import (
	"context"
	"database/sql"
	"fmt"
)

const createTableHashtagPerPostTmpl = `
CREATE TABLE IF NOT EXISTS %s
(
	post_id bigint NOT NULL,
	hashtag_id bigint NOT NULL,
	PRIMARY KEY (post_id, hashtag_id),
    CONSTRAINT fk_post
      FOREIGN KEY(post_id) 
	  REFERENCES %s(id),
    CONSTRAINT fk_hashtag
      FOREIGN KEY(hashtag_id) 
	  REFERENCES %s(id)
);
`

type StmtHashtagPerPost struct {
	stmtAddHashtagToPost *sql.Stmt
}

func (shpp *StmtHashtagPerPost) prepare(dbConn *sql.DB, hashtagPerPostTableName string) (err error) {
	const addHashtagToPost = `
		INSERT INTO %s (post_id, hashtag_id)
		VALUES ($1, $2)
`

	if shpp.stmtAddHashtagToPost, err = dbConn.Prepare(fmt.Sprintf(addHashtagToPost, hashtagPerPostTableName)); err != nil {
		return fmt.Errorf("prepare 'add hashtag to post' stmt: %w", err)
	}

	return nil
}

func (shpp *StmtHashtagPerPost) Close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = shpp.stmtAddHashtagToPost.Close(); err != nil {
			err = fmt.Errorf("close stmt 'add hashtag to post': %w", err)
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
