package log

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
}

type Formatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// NoFieldsSpace - no space between fields
	NoFieldsSpace bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool

	// CallerFirst - print caller info first
	CallerFirst bool

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string

	// ShowTimestamp - show the timestamp in the log line
	ShowTimestamp bool

	// TabbedFields - show fields in a nice format
	TabbedFields bool
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	// write time
	if f.ShowTimestamp {
		b.WriteString(entry.Time.Format(timestampFormat))
	}

	// write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if f.CallerFirst {
		f.writeCaller(b, entry)
	}

	if !f.NoColors {
		fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("]")

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors {
		b.WriteString("\x1b[0m")
	}

	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors {
		b.WriteString("\x1b[0m")
	}

	// write message
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	if !f.CallerFirst {
		f.writeCaller(b, entry)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				" (%s:%d %s)",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		if f.TabbedFields {
			b.WriteString("\n")
		}

		for _, field := range fields {

			if f.TabbedFields {
				b.WriteString("\t")
				f.writeField(b, entry, field)
				b.WriteString("\n")
			} else {
				f.writeField(b, entry, field)
			}

		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	if f.TabbedFields {
		b.WriteString("\n")
	}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			if f.TabbedFields {
				b.WriteString("\t")
				f.writeField(b, entry, field)
				b.WriteString("\n")
			} else {
				f.writeField(b, entry, field)
			}
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		if f.TabbedFields {
			b.WriteString("\n")
		}
		for _, field := range notFoundFields {
			if f.TabbedFields {
				b.WriteString("\t")
				f.writeField(b, entry, field)
				b.WriteString("\n")
			} else {
				f.writeField(b, entry, field)
			}
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {

	if field == "component" {
		fmt.Fprintf(b, "\x1b[%dm", getComponentColor(fmt.Sprintf("%v", entry.Data[field])))
	}

	if f.HideKeys {
		fmt.Fprintf(b, "[%v]", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
	}

	if field == "component" {
		b.WriteString("\x1b[0m")
	}

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getComponentColor(component string) int {
	switch component {
	case "radarr":
		return colorYellow
	case "sonarr":
		return colorBlue
	case "prowlarr":
		return colorRed
	default:
		return colorGray
	}
}

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
