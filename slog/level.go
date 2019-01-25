package slog

// Level of severity.
type Level int

const (
	_debugLevel Level = iota
	_infoLevel
	_warnLevel
	_errorLevel
	_fatalLevel
)

var levelNames = [...]string{
	_debugLevel: "DEBUG",
	_infoLevel:  "INFO",
	_warnLevel:  "WARN",
	_errorLevel: "ERROR",
	_fatalLevel: "FATAL",
}

// String implementation.
func (l Level) String() string {
	return levelNames[l]
}
