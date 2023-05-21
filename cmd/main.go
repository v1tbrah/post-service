package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/post-service/config"
	"gitlab.com/pet-pr-social-network/post-service/internal/api"
	"gitlab.com/pet-pr-social-network/post-service/internal/msgsndr"
	"gitlab.com/pet-pr-social-network/post-service/internal/storage"
)

func main() {
	newConfig := config.NewDefaultConfig()
	zerolog.SetGlobalLevel(newConfig.LogLvl)

	if err := newConfig.ParseEnv(); err != nil {
		log.Fatal().Err(err).Msg("config.ParseEnv")
	}
	zerolog.SetGlobalLevel(newConfig.LogLvl)

	ctxStart, ctxStartCancel := context.WithCancel(context.Background())

	newStorage, err := storage.Init(newConfig.StorageConfig)
	if err != nil {
		log.Fatal().Err(err).Str("config", fmt.Sprintf("%+v", newConfig.StorageConfig)).Msg("storage.Init")
	} else {
		log.Info().Msg("storage initialized")
	}

	newMsgSender, err := msgsndr.New(ctxStart, newConfig.KafkaConfig)
	if err != nil {
		log.Fatal().Err(err).Interface("config", newConfig.KafkaConfig).Msg("msgsndr.New")
	} else {
		log.Info().Msg("message sender initialized")
	}

	newAPI := api.New(newStorage, newMsgSender)

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errServingCh := make(chan error)
	go func() {
		errServing := newAPI.StartServing(ctxStart, newConfig.GRPCConfig, shutdownSig)
		errServingCh <- errServing
	}()

	select {
	case <-shutdownSig:
		close(shutdownSig)
	case errServing := <-errServingCh:
		if errServing != nil {
			log.Error().Err(errServing).Msg("newAPI.StartServing")
		}
	}

	ctxStartCancel()

	ctxClose, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err = newAPI.GracefulStop(ctxClose); err != nil {
		log.Error().Err(err).Msg("gRPC server graceful stop")
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Info().Msg("gRPC server gracefully stopped")
	}

	if err = newStorage.Close(ctxClose); err != nil {
		log.Error().Err(err).Msg("storage close")
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Info().Msg("storage closed")
	}

	if err = newMsgSender.Close(ctxClose); err != nil {
		log.Error().Err(err).Msg("msg sender close")
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Info().Msg("msg sender closed")
	}
}
