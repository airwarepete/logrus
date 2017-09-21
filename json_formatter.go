package logrus

import (
	"encoding/json"
	"fmt"
)

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}
	fields := []Field{}
	fields = append(fields, Field{Key: "time", Value: entry.Time.Format(timestampFormat)})
	fields = append(fields, Field{Key: "msg", Value: entry.Message})
	fields = append(fields, Field{Key: "level", Value: entry.Level.String()})
	for _, field := range entry.Data {
		switch v := field.Value.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			fields = append(fields, Field{Key: field.Key, Value: v.Error()})
		default:
			fields = append(fields, field)
		}
	}

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
