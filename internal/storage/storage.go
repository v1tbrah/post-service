package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gitlab.com/pet-pr-social-network/post-service/config"
)

type Storage struct {
	dbConn *sql.DB

	stmtPost           StmtPost
	stmtHashtag        StmtHashtag
	stmtHashtagPerPost StmtHashtagPerPost

	cfg config.StorageConfig
}

func Init(cfg config.StorageConfig) (*Storage, error) {
	newStorage := &Storage{cfg: cfg}

	dbConn, err := sql.Open("pgx", connString(cfg))
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	newStorage.dbConn = dbConn

	dbConn.SetMaxOpenConns(maxOpenConns)
	dbConn.SetMaxIdleConns(maxIdleConns)
	dbConn.SetConnMaxIdleTime(maxIdleTime)
	dbConn.SetConnMaxLifetime(maxLifeTime)

	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("dbConn.Ping: %w", err)
	}

	tx, err := dbConn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(fmt.Sprintf(createTablePostTmpl, cfg.PostTableName)); err != nil {
		return nil, fmt.Errorf("CreateTablePost: %w", err)
	}

	if _, err = tx.Exec(fmt.Sprintf(createTableHashtagTmpl, cfg.HashtagTableName)); err != nil {
		return nil, fmt.Errorf("CreateTableHashtag: %w", err)
	}

	if _, err = tx.Exec(fmt.Sprintf(createTableHashtagPerPostTmpl, cfg.HashtagPerPostTableName, cfg.PostTableName, cfg.HashtagTableName)); err != nil {
		return nil, fmt.Errorf("CreateTableHashtagPerPost: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if err = newStorage.stmtPost.prepare(dbConn, cfg.PostTableName); err != nil {
		return nil, fmt.Errorf("prepare 'post' stmts: %w", err)
	}

	if err = newStorage.stmtHashtag.prepare(dbConn, cfg.HashtagTableName); err != nil {
		return nil, fmt.Errorf("prepare 'hashtag' stmts: %w", err)
	}

	if err = newStorage.stmtHashtagPerPost.prepare(dbConn, cfg.HashtagPerPostTableName); err != nil {
		return nil, fmt.Errorf("prepare 'hashtag per post' stmts: %w", err)
	}

	return newStorage, nil
}

func (s *Storage) Close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = s.stmtPost.Close(ctx); err != nil {
			err = fmt.Errorf("close stmt post: %w", err)
			closeEnded <- struct{}{}
			return
		}

		if err = s.stmtHashtag.Close(ctx); err != nil {
			err = fmt.Errorf("close stmt hashtag: %w", err)
			closeEnded <- struct{}{}
			return
		}

		if err = s.stmtHashtagPerPost.Close(ctx); err != nil {
			err = fmt.Errorf("close stmt hashtag per post: %w", err)
			closeEnded <- struct{}{}
			return
		}

		if err = s.dbConn.Close(); err != nil {
			err = fmt.Errorf("close db conn: %w", err)
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

func connString(cfg config.StorageConfig) string {
	// example: "postgres://username:password@localhost:5432/database_name"
	return "postgres://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.PostDBName
}
