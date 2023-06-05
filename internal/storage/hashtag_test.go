//go:build with_db

package storage

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/pet-pr-social-network/post-service/internal/model"
)

func TestStorage_CreateHashtag(t *testing.T) {
	ctx := context.Background()

	s := tHelperInitEmptyDB(t)

	testHashtag := model.Hashtag{Name: "testHashtagName"}
	idNewHashtag, err := s.CreateHashtag(ctx, testHashtag)
	if err != nil {
		t.Fatalf("s.CreateHashtag: %v", err)
	}

	row := s.db.QueryRow(fmt.Sprintf("SELECT name FROM table_hashtag WHERE id=%d", idNewHashtag))
	if err = row.Scan(&testHashtag.Name); err != nil {
		t.Fatalf("scan get new hashtag name: %v", err)
	}
	if row.Err() != nil {
		t.Fatalf("check scan get new hashtag name: %v", row.Err())
	}

	if testHashtag.Name != "testHashtagName" {
		t.Fatalf("new hashtag name: got: %s, expected: %s", testHashtag.Name, "testHashtagName")
	}
}

func TestStorage_GetHashtag(t *testing.T) {
	ctx := context.Background()

	s := tHelperInitEmptyDB(t)

	testHashtag := model.Hashtag{Name: "testHashtagName"}
	row := s.db.QueryRow(fmt.Sprintf("INSERT INTO table_hashtag (name) VALUES('%s') RETURNING id", testHashtag.Name))
	if err := row.Scan(&testHashtag.ID); err != nil {
		t.Fatalf("scan new hashtag id: %v", err)
	}
	if row.Err() != nil {
		t.Fatalf("check scan new hashtag id: %v", row.Err())
	}

	testHashtag, err := s.GetHashtag(ctx, testHashtag.ID)
	if err != nil {
		t.Fatalf("get hashtag: %v", err)
	}

	if testHashtag.Name != "testHashtagName" {
		t.Fatalf("get hashtag name: got: %s, expected: %s", testHashtag.Name, "testHashtagName")
	}
}
