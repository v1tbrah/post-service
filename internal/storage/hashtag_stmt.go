package storage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type hashtag struct {
	create *sql.Stmt
	get    *sql.Stmt
}

func (si *hashtag) prepare(db *sql.DB) (err error) {
	const createHashtag = `
		INSERT INTO table_hashtag (name)
		VALUES ($1)
		RETURNING id;
`

	if si.create, err = db.Prepare(createHashtag); err != nil {
		return errors.Wrap(err, "prepare 'create' stmt")
	}

	const getHashtag = `
		SELECT
			id,
			name
		FROM table_hashtag
		WHERE id = $1
`

	if si.get, err = db.Prepare(getHashtag); err != nil {
		return errors.Wrap(err, "prepare 'get' stmt")
	}

	return nil
}

func (si *hashtag) close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = si.create.Close(); err != nil {
			err = errors.Wrap(err, "close stmt 'create'")
			closeEnded <- struct{}{}
			return
		}

		if err = si.get.Close(); err != nil {
			err = errors.Wrap(err, "close stmt 'get'")
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
