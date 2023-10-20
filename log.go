package serrors

import "log/slog"

func LogError(msg string, err error, args ...any) {
	args = append(args, KeyVals(err)...)
	slog.Error(msg, args...)
}
