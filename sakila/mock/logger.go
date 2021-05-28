package mock

// Logger is a mock logger.
type Logger struct {
	ErrorFn func(args ...interface{})
}

func (logger *Logger) Error(args ...interface{}) {
	if fn := logger.ErrorFn; fn != nil {
		fn(args)
	}
}
