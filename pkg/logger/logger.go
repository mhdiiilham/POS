package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type RequestID string

var RequestIDKey RequestID = "request-id"

func Info(ctx context.Context, caller, format string, values ...interface{}) {
	requestID := "-"
	ctxVal := ctx.Value(RequestIDKey)
	if ctxVal != nil {
		requestID = ctxVal.(string)
	}

	logrus.WithFields(logrus.Fields{
		"caller":     caller,
		"request_id": requestID,
	}).Infof(format, values...)
}

func Error(ctx context.Context, caller, format string, values ...interface{}) {
	requestID := "-"
	ctxVal := ctx.Value(RequestIDKey)
	if ctxVal != nil {
		requestID = ctxVal.(string)
	}

	logrus.WithFields(logrus.Fields{
		"caller":     caller,
		"request_id": requestID,
	}).Errorf(format, values...)
}

func Warn(ctx context.Context, caller, format string, values ...interface{}) {
	requestID := "-"
	ctxVal := ctx.Value(RequestIDKey)
	if ctxVal != nil {
		requestID = ctxVal.(string)
	}

	logrus.WithFields(logrus.Fields{
		"caller":     caller,
		"request_id": requestID,
	}).Warnf(format, values...)
}
