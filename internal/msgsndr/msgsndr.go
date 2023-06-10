package msgsndr

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"

	"github.com/v1tbrah/post-service/config"
)

type Sender struct {
	liveCtx context.Context

	topicPostCreatedConn *kafka.Conn
	topicPostDeletedConn *kafka.Conn

	cfg config.Kafka
}

func New(ctx context.Context, cfg config.Kafka) (*Sender, error) {
	if !cfg.Enable {
		return nil, nil
	}

	topicPostCreatedConn, err := kafka.DialLeader(ctx, "tcp", net.JoinHostPort(cfg.Host, cfg.Port), cfg.TopicPostCreated, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "kafka.DialLeader, address (%s), topic (%s)", net.JoinHostPort(cfg.Host, cfg.Port), cfg.TopicPostCreated)
	}

	topicPostDeletedConn, err := kafka.DialLeader(ctx, "tcp", net.JoinHostPort(cfg.Host, cfg.Port), cfg.TopicPostDeleted, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "kafka.DialLeader, address (%s), topic (%s)", net.JoinHostPort(cfg.Host, cfg.Port), cfg.TopicPostDeleted)
	}

	return &Sender{
		liveCtx:              ctx,
		topicPostCreatedConn: topicPostCreatedConn,
		topicPostDeletedConn: topicPostDeletedConn,
		cfg:                  cfg,
	}, nil
}

func (ms *Sender) Close(ctx context.Context) (err error) {
	if ms == nil {
		return nil
	}

	closed := make(chan struct{})

	go func() {
		if err = ms.topicPostCreatedConn.Close(); err != nil {
			err = errors.Wrap(err, "topicPostCreatedConn.close")
			closed <- struct{}{}
			return
		}

		if err = ms.topicPostDeletedConn.Close(); err != nil {
			err = errors.Wrap(err, "topicPostDeletedConn.close")
			closed <- struct{}{}
			return
		}

		closed <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-closed:
		return err
	}
}
