package def

type ServerType int

const (
	Client ServerType = iota
	Gateway
	Hub
)
