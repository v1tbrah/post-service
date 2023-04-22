package storage

import (
	"context"
	"database/sql"
	"fmt"
)

const createTableHashtagTmpl = `
CREATE TABLE IF NOT EXISTS %s
	(
		id   serial  PRIMARY KEY,
		name varchar UNIQUE NOT NULL
	);
`

type StmtHashtag struct {
	stmtCreateHashtag *sql.Stmt
	stmtGetHashtag    *sql.Stmt
}

func (si *StmtHashtag) prepare(dbConn *sql.DB, hashtagTableName string) (err error) {
	const createHashtag = `
		INSERT INTO %s (name)
		VALUES ($1)
		RETURNING id;
`

	if si.stmtCreateHashtag, err = dbConn.Prepare(fmt.Sprintf(createHashtag, hashtagTableName)); err != nil {
		return fmt.Errorf("prepare 'create hashtag' stmt: %w", err)
	}

	const getHashtag = `
		SELECT
			hashtag.id,
			hashtag.name
		FROM %s AS hashtag
		WHERE hashtag.id = $1
`

	if si.stmtGetHashtag, err = dbConn.Prepare(fmt.Sprintf(getHashtag, hashtagTableName)); err != nil {
		return fmt.Errorf("prepare 'get hashtag' stmt: %w", err)
	}

	return nil
}

func (si *StmtHashtag) Close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = si.stmtCreateHashtag.Close(); err != nil {
			err = fmt.Errorf("close stmt 'create hashtag': %w", err)
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
