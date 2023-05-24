package types

type Status int64

const (
	Unknown Status = iota
	OK
	Warning
	Critical
)
