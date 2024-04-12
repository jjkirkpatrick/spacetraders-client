package client

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/phuslu/log"
)

func Logformat(out io.Writer, args *log.FormatterArgs) (n int, err error) {

	const (
		Reset   = "\x1b[0m"
		Black   = "\x1b[30m"
		Red     = "\x1b[31m"
		Green   = "\x1b[32m"
		Yellow  = "\x1b[33m"
		Blue    = "\x1b[34m"
		Magenta = "\x1b[35m"
		Cyan    = "\x1b[36m"
		White   = "\x1b[37m"
		Gray    = "\x1b[90m"
	)

	// colorful level string
	var color, three string
	switch args.Level {
	case "trace":
		color, three = Magenta, "TRACE"
	case "debug":
		color, three = Yellow, "DEBUG"
	case "info":
		color, three = Green, "INFO"
	case "warn":
		color, three = Red, "WARN"
	case "error":
		color, three = Red, "ERROR"
	case "fatal":
		color, three = Red, "FATAL"
	case "panic":
		color, three = Red, "PANIC"
	default:
		color, three = Gray, "???"
	}

	b := &strings.Builder{}
	// pretty console writer
	// header
	fmt.Fprintf(b, "%s%s%s %s%s%s ", Gray, args.Time, Reset, color, three, Reset)
	if args.Caller != "" {
		fmt.Fprintf(b, "%s %s %s>%s", args.Goid, args.Caller, Cyan, Reset)
	} else {
		fmt.Fprintf(b, "%s>%s", Cyan, Reset)
	}

	// message
	fmt.Fprintf(b, " %s", args.Message)

	// key and values
	for _, kv := range args.KeyValues {
		if kv.ValueType == 's' {
			kv.Value = strconv.Quote(kv.Value)
		}
		if kv.Key == "error" {
			fmt.Fprintf(b, " %s%s=%s%s", Red, kv.Key, kv.Value, Reset)
		} else {
			fmt.Fprintf(b, " %s%s=%s%s%s", Cyan, kv.Key, Gray, kv.Value, Reset)
		}
	}

	// stack
	if args.Stack != "" {
		b.WriteString("\n")
		b.WriteString(args.Stack)
		if args.Stack[len(args.Stack)-1] != '\n' {
			b.WriteString("\n")
		}
	} else {
		b.WriteString("\n")
	}
	return out.Write([]byte(b.String()))
}
