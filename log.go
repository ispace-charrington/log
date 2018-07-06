package log

import "os"
import "fmt"
import "time"

// A Logger captures program events at varying severity levels, and
// is relatively simple to nest to indicate logic structure.
type Logger interface {
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Panic(format string, v ...interface{})
	Prefix(prefix string) Logger
}

// DefaultLogger is a simple logger that discards Info/Warning/Panic
// metadata (aside from panicing on Panic() of course) and simply
// writes the timestamped log data to stderr.
type DefaultLogger struct {
	prefix string
}

// Default returns a logger suitable for writing to stderr.
func Default() (l *DefaultLogger) {
	l = &DefaultLogger{}
	return
}

func (l *DefaultLogger) out(t string) {
	_, err := fmt.Fprintf(
		os.Stderr,
		"%s\t%s\t%s\n",
		time.Now().UTC().Format(time.RFC3339),
		l.prefix, t)

	if err != nil {
		panic(fmt.Sprintf("Failed to log to stderr!\nError: %v\nLog: %s\n", err, t))
	}
}

// Info writes to os.Stderr, but does not record anything more or
// less important than other logging levels.
func (l *DefaultLogger) Info(format string, v ...interface{}) {
	l.out(fmt.Sprintf(format, v...))
}

// Warning writes to os.Stderr, but does not record anything more or
// less important than other logging levels.
func (l *DefaultLogger) Warning(format string, v ...interface{}) {
	l.out(fmt.Sprintf(format, v...))
}

// Panic writes to os.Stderr, and then panics. It doesn't record
// anything more or less important than other logging levels. But
// it does panic, so steady now.
func (l *DefaultLogger) Panic(format string, v ...interface{}) {
	t := fmt.Sprintf(format, v...)
	l.out(t)
	// if stderr is redirected, let's flush to storage in case we
	// are about to reboot or crash
	os.Stderr.Sync()
	panic(t)
}

// Prefix returns a new DefaultLogger with this prefix appended.
func (l *DefaultLogger) Prefix(prefix string) *DefaultLogger {
	if l.prefix == "" {
		return &DefaultLogger{prefix: prefix}
	}
	return &DefaultLogger{prefix: fmt.Sprintf("%s:%s", l.prefix, prefix)}
}

// A NullLogger discards all Info and Warning logs, and simply
// panics all Panic logs.
type NullLogger struct{}

// Null returns a fully-initialized, ready-to-use, standards
// compliant, community endorsed, efficient, scalable, redundant,
// HTML 4.1 Transitional DOCTYPE declared NullLogger that discards
// all log data.
func Null() *NullLogger {
	return &NullLogger{}
}

// Info discards the logged data.
func (n *NullLogger) Info(format string, v ...interface{}) {}

// Warning discards the logged data.
func (n *NullLogger) Warning(format string, v ...interface{}) {}

// Panic formats the data into a string, then panics.
func (n *NullLogger) Panic(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}

// Prefix returns a pointer to the NullLogger and discards the
// provided prefix.
func (n *NullLogger) Prefix(prefix string) *NullLogger {
	return n
}
