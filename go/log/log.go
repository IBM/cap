/*
Copyright 2018 The cap Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations
*/

package log

import (
	"context"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

var (
	// G is an alias for GetLogger.
	G = GetLogger

	// L is an alias for the the standard logger.
	L = logrus.NewEntry(logrus.StandardLogger())
)

type (
	loggerKey struct{}
)

// TraceLevel - trick to overload logrus default levels to add trace
const TraceLevel = logrus.Level(uint32(logrus.DebugLevel + 1))

// ParseLevel takes a string level and returns the Logrus log level constant.
// It supports trace level.
func ParseLevel(lvl string) (logrus.Level, error) {
	if lvl == "trace" {
		return TraceLevel, nil
	}
	return logrus.ParseLevel(lvl)
}

// PrintLevel takes a logrus.Logger and returns the Logrus log level as a string
// It supports trace level.
func PrintLevel(l *logrus.Logger) string {
	if l.Level == TraceLevel {
		return "trace"
	}
	return l.Level.String()
}

// WithLogger returns a new context with the provided logger. To be used in
// combination with logger.WithField(s).
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		return L
	}

	return logger.(*logrus.Entry)
}

// Trace logs a message at level Trace with the log entry passed-in.
func Trace(e *logrus.Entry, args ...interface{}) {
	level := logrus.Level(atomic.LoadUint32((*uint32)(&e.Logger.Level)))
	if level >= TraceLevel {
		e.Debug(args...)
	}
}

// Tracef logs a message at level Trace with the log entry passed-in.
func Tracef(e *logrus.Entry, format string, args ...interface{}) {
	level := logrus.Level(atomic.LoadUint32((*uint32)(&e.Logger.Level)))
	if level >= TraceLevel {
		e.Debugf(format, args...)
	}
}
