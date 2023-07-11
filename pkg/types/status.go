package types

import "github.com/gitu/katastasi/pkg/config"

type Status int64

const (
	Unknown Status = iota
	OK
	Info
	Warning
	Critical
)

func MapStatus(severity config.Severity) Status {
	switch severity {
	case config.Critical:
		return Critical
	case config.Warning:
		return Warning
	case config.Info:
		return Info
	case config.OK:
		return OK
	case config.Unknown:
		return Unknown
	}
	return Unknown
}
