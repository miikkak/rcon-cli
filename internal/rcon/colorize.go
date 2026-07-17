package rcon

import (
	"runtime"
	"strconv"
	"strings"
)

// sectionSign is Minecraft's formatting-code prefix character.
const sectionSign = '§'

// reset is the ANSI escape that clears all formatting.
const reset = "\x1b[0m"

// colors maps Minecraft's legacy single-character formatting codes to their
// ANSI escape equivalents.
var colors = map[rune]string{
	'0': "\x1b[30m", // black
	'1': "\x1b[34m", // dark blue
	'2': "\x1b[32m", // dark green
	'3': "\x1b[36m", // dark aqua
	'4': "\x1b[31m", // dark red
	'5': "\x1b[35m", // dark purple
	'6': "\x1b[33m", // gold
	'7': "\x1b[37m", // gray
	'8': "\x1b[30m", // dark gray
	'9': "\x1b[34m", // blue
	'a': "\x1b[32m", // green
	'b': "\x1b[32m", // aqua
	'c': "\x1b[31m", // red
	'd': "\x1b[35m", // light purple
	'e': "\x1b[33m", // yellow
	'f': "\x1b[37m", // white
	'k': "",         // random
	'm': "\x1b[9m",  // strikethrough
	'o': "\x1b[3m",  // italic
	'l': "\x1b[1m",  // bold
	'n': "\x1b[4m",  // underline
	'r': reset,      // reset
}

// colorize translates Minecraft formatting codes embedded in str into ANSI
// escape sequences.
//
// It recognizes two kinds of codes:
//   - §x§R§R§G§G§B§B (truecolor: §x followed by six §<hex digit> pairs,
//     used by e.g. BlueMap) -> \x1b[38;2;R;G;Bm
//   - §<legacy code> (one of the keys in colors) -> the mapped ANSI escape
//
// A §x not followed by a complete, well-formed truecolor sequence is left
// untouched rather than partially consumed, so it never corrupts nearby
// legacy codes or panics on truncated input. An unrecognized §<char> is
// likewise left untouched, matching the historical behavior of only
// substituting known codes.
func colorize(str string) string {
	if runtime.GOOS == "windows" {
		return str
	}

	runes := []rune(str)
	var out strings.Builder
	out.Grow(len(str))

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r != sectionSign {
			out.WriteRune(r)
			continue
		}

		if i+1 >= len(runes) {
			out.WriteRune(r)
			continue
		}

		next := runes[i+1]
		if next == 'x' || next == 'X' {
			if hex, ok := parseTrueColor(runes, i); ok {
				out.WriteString(trueColorEscape(hex))
				i += 13 // §x + 6×(§h) = 14 runes total; loop's i++ covers the 14th
				continue
			}
			out.WriteRune(r)
			continue
		}

		if ansi, ok := colors[next]; ok {
			out.WriteString(ansi)
			i++
			continue
		}

		out.WriteRune(r)
	}

	return strings.ReplaceAll(out.String(), "\n", "\n"+reset)
}

// parseTrueColor checks runes[start:] for §x§h§h§h§h§h§h (start points at
// the § before 'x'). It returns the six hex digits concatenated (RRGGBB)
// and whether the sequence was complete and well-formed.
func parseTrueColor(runes []rune, start int) (hex string, ok bool) {
	const numPairs = 6
	need := 2 + numPairs*2 // §x + 6×(§h)
	if start+need > len(runes) {
		return "", false
	}

	var b strings.Builder
	pos := start + 2 // skip §x
	for range numPairs {
		if runes[pos] != sectionSign {
			return "", false
		}
		h := runes[pos+1]
		if !isHexDigit(h) {
			return "", false
		}
		b.WriteRune(h)
		pos += 2
	}

	return b.String(), true
}

func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func trueColorEscape(hex string) string {
	rr, _ := strconv.ParseUint(hex[0:2], 16, 8)
	gg, _ := strconv.ParseUint(hex[2:4], 16, 8)
	bb, _ := strconv.ParseUint(hex[4:6], 16, 8)
	return "\x1b[38;2;" + strconv.FormatUint(rr, 10) + ";" +
		strconv.FormatUint(gg, 10) + ";" + strconv.FormatUint(bb, 10) + "m"
}
