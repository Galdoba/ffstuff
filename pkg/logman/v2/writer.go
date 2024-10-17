package logman

import "io"

type writerInfo struct {
	writer             io.Writer
	writeFieldsRequest []string
	fieldColorMap      map[string]bool
	fieldFormatterMap  map[string]func(string, ...any) string
	outputFields       []string
}
