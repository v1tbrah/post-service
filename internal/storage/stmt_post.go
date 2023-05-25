package storage

import (
	"context"
	"database/sql"
	"fmt"
)

const createTablePostTmpl = `
CREATE TABLE IF NOT EXISTS %s
	(
		id serial PRIMARY KEY,
		user_id bigint NOT NULL,
		description varchar,
		created_at timestamp NOT NULL
	);
`

type StmtPost struct {
	stmtCreatePost       *sql.Stmt
	stmtGetPost          *sql.Stmt
	stmtGetPostsByUserID *sql.Stmt
}

func (sc *StmtPost) prepare(dbConn *sql.DB, postTableName string) (err error) {
	const createPost = `
		INSERT INTO %s (user_id, description, created_at)
		VALUES ($1, $2, $3)
		RETURNING id;
`

	if sc.stmtCreatePost, err = dbConn.Prepare(fmt.Sprintf(createPost, postTableName)); err != nil {
		return fmt.Errorf("prepare 'create post' stmt: %w", err)
	}

	const getPost = `
		SELECT
			post.id,
			post.user_id,
			post.description,
			post.created_at
		FROM %s AS post
		WHERE post.id = $1
`

	if sc.stmtGetPost, err = dbConn.Prepare(fmt.Sprintf(getPost, postTableName)); err != nil {
		return fmt.Errorf("prepare 'get post' stmt: %w", err)
	}

	const getPostsByUserID = `
		SELECT
			post.id,
			post.user_id,
			post.description,
			post.created_at
		FROM %s AS post
		WHERE post.user_id = $1
`

	if sc.stmtGetPostsByUserID, err = dbConn.Prepare(fmt.Sprintf(getPostsByUserID, postTableName)); err != nil {
		return fmt.Errorf("prepare 'get posts by user id' stmt: %w", err)
	}

	return nil
}

func (sc *StmtPost) Close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = sc.stmtCreatePost.Close(); err != nil {
			err = fmt.Errorf("close stmt 'create post': %w", err)
			closeEnded <- struct{}{}
			return
		}

		if err = sc.stmtGetPost.Close(); err != nil {
			err = fmt.Errorf("close stmt 'get post': %w", err)
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
