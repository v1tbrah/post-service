package msgsndr

import (
	"encoding/binary"
	"time"

	"github.com/avast/retry-go"
	"github.com/rs/zerolog/log"
)

const (
	postDeleteRetryDelay    = time.Second * 10
	postDeleteRetryAttempts = 5
)

func (ms *Sender) SendMsgPostDeleted(id int64) {
	if ms == nil {
		return
	}

	msgBody := make([]byte, 8)
	binary.PutVarint(msgBody, id)

	err := retry.Do(func() error {
		_, err := ms.topicPostDeletedConn.Write(msgBody)
		return err
	},
		retry.Context(ms.liveCtx),
		retry.Delay(postDeleteRetryDelay),
		retry.Attempts(postDeleteRetryAttempts),
	)

	if err != nil {
		log.Warn().Err(err).Int64("id", id).Msg("topicPostDeletedConn.Write")
		return
	}
}
