// Copyright (c) 2013 Raffaele Sena https://github.com/raff
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish, dis-
// tribute, sublicense, and/or sell copies of the Software, and to permit
// persons to whom the Software is furnished to do so, subject to the fol-
// lowing conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABIL-
// ITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT
// SHALL THE AUTHOR BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/ansi"
)

//
// Concatenate pattern values and build
// a single regexp to find any of them (pattern1 or pattern2 or ...)
//
func makePatternLevels(patternMap map[string]string) *regexp.Regexp {
	keys := make([]string, 0, len(patternMap))
	for key, _ := range patternMap {
		keys = append(keys, key)
	}

	return regexp.MustCompile(" (" + strings.Join(keys, "|") + ") ")
}

var (
	color_custom   = ansi.ColorCode("white+h:blue")
	color_info     = ansi.ColorCode("green+h:black")
	color_warn     = ansi.ColorCode("yellow+h:black")
	color_error    = ansi.ColorCode("red+h:black")
	color_critical = ansi.ColorCode("yellow+h:red")
	color_fatal    = ansi.ColorCode("orange+h:red")
	color_debug    = ansi.ColorCode("cyan+h:black")
	color_trace    = ansi.ColorCode("blue+h:black")
	color_reset    = ansi.ColorCode("reset")
)

var (
	levels = map[string]string{
		"INFO":     color_info,
		"WARN":     color_warn,
		"WARNING":  color_warn,
		"ERROR":    color_error,
		"CRITICAL": color_critical,
		"FATAL":    color_fatal,
		"DEBUG":    color_debug,
		"TRACE":    color_trace,
	}

	pattern_level  = makePatternLevels(levels)
	pattern_custom *regexp.Regexp
)

func Colorize(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if pattern_custom != nil {
			if match := pattern_custom.FindString(line); len(match) > 0 {
				fmt.Println(color_custom + line + color_reset)
				continue
			}
		}

		if match := pattern_level.FindString(line); len(match) > 0 {
			if color := levels[strings.TrimSpace(match)]; len(color) > 0 {
				fmt.Println(color + line + color_reset)
				continue
			}
		}

		fmt.Println(line)
	}
}

func main() {

	if len(os.Args) > 1 {
		pattern_custom = regexp.MustCompile(os.Args[1])
	}

	Colorize(os.Stdin)
}
