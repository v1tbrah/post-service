//go:build with_db

package storage

import (
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

	s, err := Init(cfg.Storage)
	if err != nil {
		t.Fatalf("init storage: %v", err)
	}

	// CHECK EXISTENCE AFTER INITIALIZATION
	if _, err = s.db.Query("SELECT 1 FROM table_post"); err != nil {
		t.Fatalf("select from table post: %v", err)
	}

	if _, err = s.db.Query("SELECT 1 FROM table_hashtag"); err != nil {
		t.Fatalf("select from table hashtag: %v", err)
	}

	if _, err = s.db.Query("SELECT 1 FROM table_hashtag_per_post"); err != nil {
		t.Fatalf("select from table hashtag per post: %v", err)
	}
}

func tHelperInitEmptyDB(t *testing.T) *Storage {
	cfg := config.NewDefaultConfig()
	zerolog.SetGlobalLevel(cfg.LogLvl)

	if err := cfg.ParseEnv(); err != nil {
		log.Fatal().Err(err).Msg("config.ParseEnv")
	}
	zerolog.SetGlobalLevel(cfg.LogLvl)

	s, err := Init(cfg.Storage)
	if err != nil {
		t.Fatalf("init storage: %v", err)
	}

	// DELETE FROM TABLES FOR CLEAR TEST SPACE
	if _, err = s.db.Query("DELETE FROM table_hashtag_per_post CASCADE"); err != nil {
		t.Fatalf("delete from table hashtag per post: %v", err)
	}

	if _, err = s.db.Query("DELETE FROM table_post CASCADE"); err != nil {
		t.Fatalf("delete from table post: %v", err)
	}

	if _, err = s.db.Query("DELETE FROM table_hashtag CASCADE"); err != nil {
		t.Fatalf("delete from table hashtag: %v", err)
	}

	return s
}
