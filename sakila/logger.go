package sakila

// Logger is a sakila logger.
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}
