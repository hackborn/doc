package doc

import (
	"fmt"
	"strings"
)

type StringifyArgs struct {
	Separator string
	Quote     rune
}

var (
	emptyStringifyArgs = StringifyArgs{}
)

// writeRune writes the supplied rune to the string builder,
// answering err for any error. It will skip the rune if it's
// emptyStringifyArgs.Quote.
func writeRune(err error, r rune, sb *strings.Builder) error {
	if err != nil {
		return err
	}
	if r == emptyStringifyArgs.Quote {
		return nil
	}
	_, err = sb.WriteRune(r)
	return err
}
func writeString(err error, s string, sb *strings.Builder) error {
	if err != nil {
		return err
	}
	_, err = sb.WriteString(s)
	return err
}

// ExpandPairs compiles all values in pairs into a single
// comma-separated string. There must be an even number of
// pairs, which are treated as key=values. sep is
// inserted between each pair. String values are quoted with quoteRune.
// (use rune(0) if you don't want quotes).
func ExpandPairs(args StringifyArgs, sb *strings.Builder, pairs ...any) string {
	onKey := true
	for i, s := range pairs {
		if i > 0 && i%2 == 0 {
			if args.Separator != emptyStringifyArgs.Separator {
				sb.WriteString(args.Separator)
			}
		}
		switch v := s.(type) {
		case string:
			if !onKey && args.Quote != emptyStringifyArgs.Quote {
				sb.WriteRune(args.Quote)
			}
			sb.WriteString(v)
			if !onKey && args.Quote != emptyStringifyArgs.Quote {
				sb.WriteRune(args.Quote)
			}
		default:
			sb.WriteString(fmt.Sprintf("%v", v))
		}
		if onKey {
			sb.WriteString("=")
		}
		onKey = !onKey
	}
	return sb.String()
}
