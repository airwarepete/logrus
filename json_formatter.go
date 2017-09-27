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
	// this translates the field key/values back to a map to make the json output more readable even tho the order may still be jacked.
	fields := make(map[string]string)
	fields["time"] = entry.Time.Format(timestampFormat)
	fields["level"] = entry.Level.String()
	fields["msg"] = entry.Message
	for _, field := range entry.Data {
		switch v := field.Value.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			fields[field.Key] = v.Error()
		default:
			fields[field.Key] = field.Value.(string)
		}
	}

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
