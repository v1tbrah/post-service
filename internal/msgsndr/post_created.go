package msgsndr

import (
	"encoding/json"
	"time"

	"github.com/avast/retry-go"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/post-service/internal/model"
)

const (
	postCreateRetryDelay    = time.Second * 10
	postCreateRetryAttempts = 5
)

func (ms *Sender) SendMsgPostCreated(post model.Post) {
	if ms == nil {
		return
	}

	msgBody, err := json.Marshal(post)
	if err != nil {
		log.Warn().Err(err).Interface("post", post).Msg("json.Marshal, post")
		return
	}

	err = retry.Do(func() error {
		_, err = ms.topicPostCreatedConn.Write(msgBody)
		return err
	},
		retry.Context(ms.liveCtx),
		retry.Delay(postCreateRetryDelay),
		retry.Attempts(postCreateRetryAttempts),
	)

	if err != nil {
		log.Warn().Err(err).Interface("post", post).Msg("topicPostCreatedConn.Write")
		return
	}
}
