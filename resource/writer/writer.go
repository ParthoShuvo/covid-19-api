package writer

import "net/http"

// WriterType defines all writers
type WriterType string

// available writer types
const (
	JSON WriterType = "json"
)

// Writer provides http-data writing operation
type Writer interface {
	Write(rw http.ResponseWriter, data interface{})
}

type newWriter func() Writer

//NewWriter is a factory function that provides writer by its type
func NewWriter(wt WriterType) Writer {
	writerMap := map[WriterType]newWriter{
		JSON: newJSONWriter,
	}
	newWriter := writerMap[wt]
	return newWriter()
}
