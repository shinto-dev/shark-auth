package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

func FromContext(ctx context.Context) *logrus.Entry {
	fields := getLogFields(ctx)
	return logrus.WithFields(fields)
}