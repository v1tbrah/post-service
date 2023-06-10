package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"

	"github.com/v1tbrah/post-service/config"
)

type Storage struct {
	db *sql.DB

	post           post
	hashtag        hashtag
	hashtagPerPost hashtagPerPost

	cfg config.Storage
}

func Init(cfg config.Storage) (*Storage, error) {
	newStorage := &Storage{cfg: cfg}

	db, err := sql.Open("pgx", connString(cfg))
	if err != nil {
		return nil, errors.Wrapf(err, "sql.Open, conn string: %s", connString(cfg))
	}
	newStorage.db = db

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)
	db.SetConnMaxLifetime(maxLifeTime)

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "db.Ping")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer tx.Rollback()

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "tx.Commit")
	}

	if err = newStorage.post.prepare(db); err != nil {
		return nil, errors.Wrap(err, "prepare 'post' stmts")
	}

	if err = newStorage.hashtag.prepare(db); err != nil {
		return nil, errors.Wrap(err, "prepare 'hashtag' stmts")
	}

	if err = newStorage.hashtagPerPost.prepare(db); err != nil {
		return nil, errors.Wrap(err, "prepare 'hashtag per post' stmts")
	}

	return newStorage, nil
}

func (s *Storage) Close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = s.post.close(ctx); err != nil {
			err = errors.Wrap(err, "close stmt 'post'")
			closeEnded <- struct{}{}
			return
		}

		if err = s.hashtag.close(ctx); err != nil {
			err = errors.Wrap(err, "close stmt 'hashtag'")
			closeEnded <- struct{}{}
			return
		}

		if err = s.hashtagPerPost.close(ctx); err != nil {
			err = errors.Wrap(err, "close stmt 'hashtag per post'")
			closeEnded <- struct{}{}
			return
		}

		if err = s.db.Close(); err != nil {
			err = errors.Wrap(err, "close db conn")
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

func connString(cfg config.Storage) string {
	// example: "postgres://username:password@localhost:5432/database_name"
	return "postgres://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.PostDBName
}
