package logman

import (
	"fmt"
	"strings"
)

// This is a convinience function for ProcessMessage.
// Printf formats message according to a format specifier and writes to output writers of Level INFO.
// It returns message processing error encountered.
// Same as func: Info(format, args...) (error) .
func Printf(format string, args ...interface{}) error {
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[INFO]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Println writes to output writers of Level INFO. Args are separated
// with spaces.
// It returns message processing error encountered.
func Println(args ...interface{}) error {
	format := ""
	for range args {
		format += "%v "
	}
	format = strings.TrimSuffix(format, " ")
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[INFO]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Fatalf formats message according to a format specifier and writes to output writers of Level FATAL.
// It returns message processing error encountered or error created if processing is success.
// By default calling level Fatal cause os.Exit(1) after completion (subject to change during logger setup process).
func Fatalf(format string, args ...interface{}) error {
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[FATAL]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Errorf formats message according to a format specifier and writes to output writers of Level ERROR.
// It returns message processing error encountered or error created if processing is success.
func Errorf(format string, args ...interface{}) error {
	errCreated := fmt.Errorf(format, args...)
	msg := NewMessage(format, args...)
	if errProcessing := process(msg, logMan.logLevels[ERROR]); errProcessing != nil {
		return errProcessing
	}
	return errCreated
}

// This is a convinience function for ProcessMessage.
// Error creates message input argument and writes to output writers of Level ERROR.
// It returns message processing error encountered or input error if processing is success.
func Error(errInput error) error {
	msg := NewMessage(errInput.Error())
	if errProcessing := process(msg, logMan.logLevels[ERROR]); errProcessing != nil {
		return errProcessing
	}
	return errInput
}

// This is a convinience function for ProcessMessage.
// Warn formats message according to a format specifier and writes to output writers of Level WARN.
// It returns message processing error encountered.
func Warn(format string, args ...interface{}) error {
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[WARN]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Info formats message according to a format specifier and writes to output writers of Level INFO.
// It returns message processing error encountered.
func Info(format string, args ...interface{}) error {
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[INFO]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Debug receives message with map of additional values and writes to output writers of Level DEBUG.
// It returns message processing error encountered.
func Debug(msg Message, values map[string]interface{}) error {
	for k, val := range values {
		msg.SetField(k, val)
	}
	if err := process(msg, logMan.logLevels[DEBUG]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Trace receives message with map of additional values and writes to output writers of Level TRACE.
// It returns message processing error encountered.
func Trace(msg Message, values map[string]interface{}) error {
	for k, val := range values {
		msg.SetField(k, val)
	}
	if err := process(msg, logMan.logLevels[TRACE]); err != nil {
		return err
	}
	return nil
}
