package logger

type DummyLogger struct{}

func (d DummyLogger) Debug(msg string, args ...any) {}
func (d DummyLogger) Info(msg string, args ...any)  {}
func (d DummyLogger) Warn(msg string, args ...any)  {}
func (d DummyLogger) Error(msg string, args ...any) {}
func (d DummyLogger) With(args ...any) Logger       { return d }
