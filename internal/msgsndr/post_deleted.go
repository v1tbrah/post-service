package msgsndr

import (
	"encoding/json"
	"time"

	"github.com/avast/retry-go"
	"github.com/rs/zerolog/log"
)

const (
	postDeleteRetryDelay    = time.Second * 10
	postDeleteRetryAttempts = 5
)

type MsgPostDeleted struct {
	ID     int64
	UserID int64
}

func (ms *Sender) SendMsgPostDeleted(id int64, userID int64) {
	if ms == nil {
		return
	}

	msg := MsgPostDeleted{ID: id, UserID: userID}

	msgBody, err := json.Marshal(msg)
	if err != nil {
		log.Warn().Err(err).Interface("msg", msg).Msg("json.Marshal, msg")
		return
	}

	err = retry.Do(func() error {
		_, err = ms.topicPostDeletedConn.Write(msgBody)
		return err
	},
		retry.Context(ms.liveCtx),
		retry.Delay(postDeleteRetryDelay),
		retry.Attempts(postDeleteRetryAttempts),
	)

	if err != nil {
		log.Warn().Err(err).Interface("msg", msg).Msg("topicPostDeletedConn.Write")
		return
	}
}
