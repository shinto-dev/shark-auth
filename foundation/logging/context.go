package logging

import (
	"context"
)

type logFieldsKeyType string

const logFieldsKey = logFieldsKeyType("key")

func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logFieldsKey, map[string]interface{}{})
}

//Set This adds new log field to the context
func Set(ctx context.Context, key string, value string) {
	fields, ok := ctx.Value(logFieldsKey).(map[string]interface{})
	if !ok {
		return
	}
	fields[key] = value
}

func getLogFields(ctx context.Context) map[string]interface{} {
	logFields := ctx.Value(logFieldsKey)
	if logFields == nil {
		return map[string]interface{}{}
	}

	return logFields.(map[string]interface{})
}
