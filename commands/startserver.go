package commands

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"shark-auth/apis"
)

func newStartServerCommand(db *sqlx.DB, redisClient *redis.Client) *cobra.Command {
	return &cobra.Command{
		Use:     "startserver",
		Short:   "Start HTTP API server",
		Aliases: []string{"startapp", "runserver"},
		Run: func(_ *cobra.Command, _ []string) {
			router := apis.API(db, redisClient)

			server := http.Server{
				Addr:    fmt.Sprintf(":%d", 8080), // todo move to config
				Handler: router,
			}

			go func() {
				logrus.Infof("starting HTTP server, listening at %d", 8080)
				if err := server.ListenAndServe(); err != nil {
					logrus.Error("failed to start the server")
				}

			}()

			sigquit := make(chan os.Signal, 1)
			signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

			_ = <-sigquit
			logrus.Info("gracefully shutting down the server")

			if err := server.Shutdown(context.Background()); err != nil {
				logrus.WithError(err).Error("unable to shutdown the server")
				return
			}

			logrus.Info("server stopped")
		},
	}
}
