//go:build with_db

package storage

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/pet-pr-social-network/post-service/config"
)

func TestStorage_Init(t *testing.T) {
	cfg := config.NewDefaultConfig()
	zerolog.SetGlobalLevel(cfg.LogLvl)

	if err := cfg.ParseEnv(); err != nil {
		log.Fatal().Err(err).Msg("config.ParseEnv")
	}
	zerolog.SetGlobalLevel(cfg.LogLvl)

	s := tHelperInitEmptyDB(t)

	// DROP TABLES TO CHECK THEIR EXISTENCE AFTER REINITIALIZATION
	if _, err := s.dbConn.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", cfg.StorageConfig.PostTableName)); err != nil {
		t.Fatalf("drop table post: %s", err)
	}

	if _, err := s.dbConn.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", cfg.StorageConfig.HashtagTableName)); err != nil {
		t.Fatalf("drop table hashtag: %s", err)
	}

	if _, err := s.dbConn.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", cfg.StorageConfig.HashtagPerPostTableName)); err != nil {
		t.Fatalf("drop table hashtag per post: %s", err)
	}
}

func tHelperInitEmptyDB(t *testing.T) *Storage {
	cfg := config.NewDefaultConfig()
	zerolog.SetGlobalLevel(cfg.LogLvl)

	if err := cfg.ParseEnv(); err != nil {
		log.Fatal().Err(err).Msg("config.ParseEnv")
	}
	zerolog.SetGlobalLevel(cfg.LogLvl)

	s, err := Init(cfg.StorageConfig)
	if err != nil {
		t.Fatalf("init storage: %v", err)
	}

	// DROP TABLES IF THEY ALREADY EXIST
	if _, err = s.dbConn.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", s.cfg.PostTableName)); err != nil {
		t.Fatalf("drop table post: %s", err)
	}

	if _, err = s.dbConn.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", s.cfg.HashtagTableName)); err != nil {
		t.Fatalf("drop table hashtag: %s", err)
	}

	if _, err = s.dbConn.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", s.cfg.HashtagPerPostTableName)); err != nil {
		t.Fatalf("drop table hashtag per post: %s", err)
	}

	// REINIT
	if s, err = Init(cfg.StorageConfig); err != nil {
		t.Fatalf("init storage after drop tables: %v", err)
	}

	return s
}
