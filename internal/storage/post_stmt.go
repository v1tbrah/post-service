package storage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type post struct {
	create          *sql.Stmt
	delete          *sql.Stmt
	get             *sql.Stmt
	getListByUserID *sql.Stmt
}

func (sc *post) prepare(dbConn *sql.DB) (err error) {
	const createPost = `
		INSERT INTO table_post (user_id, description, created_at)
		VALUES ($1, $2, $3)
		RETURNING id;
`

	if sc.create, err = dbConn.Prepare(createPost); err != nil {
		return errors.Wrap(err, "prepare 'create' stmt")
	}

	const deletePost = `DELETE FROM table_post WHERE id = $1 RETURNING user_id;`

	if sc.delete, err = dbConn.Prepare(deletePost); err != nil {
		return errors.Wrap(err, "prepare 'delete' stmt")
	}

	const getPost = `
		SELECT
			id,
			user_id,
			description,
			created_at
		FROM table_post
		WHERE id = $1;
`

	if sc.get, err = dbConn.Prepare(getPost); err != nil {
		return errors.Wrap(err, "prepare 'get' stmt")
	}

	const getPostsByUserID = `
		SELECT
			id,
			user_id,
			description,
			created_at
		FROM table_post AS post
		WHERE user_id = $1
`

	if sc.getListByUserID, err = dbConn.Prepare(getPostsByUserID); err != nil {
		return errors.Wrap(err, "prepare 'get list by user id' stmt")
	}

	return nil
}

func (sc *post) close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = sc.create.Close(); err != nil {
			err = errors.Wrap(err, "close stmt 'create post'")
			closeEnded <- struct{}{}
			return
		}

		if err = sc.get.Close(); err != nil {
			err = errors.Wrap(err, "close stmt 'get post'")
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
