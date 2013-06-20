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
        "flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/ansi"
)

type colorfn func(string) string

//
// Concatenate pattern values and build
// a single regexp to find any of them (pattern1 or pattern2 or ...)
//
func makePatternLevels(patternMap map[string]colorfn) *regexp.Regexp {
	keys := make([]string, 0, len(patternMap))
	for key, _ := range patternMap {
		keys = append(keys, key)
	}

	return regexp.MustCompile(" (" + strings.Join(keys, "|") + ") ")
}

var (
	color_custom   = ansi.ColorFunc("white+h:blue")
	color_info     = ansi.ColorFunc("green+h:black")
	color_warn     = ansi.ColorFunc("yellow+h:black")
	color_error    = ansi.ColorFunc("red+h:black")
	color_critical = ansi.ColorFunc("yellow+h:red")
	color_fatal    = ansi.ColorFunc("orange+h:red")
	color_debug    = ansi.ColorFunc("cyan+h:black")
	color_trace    = ansi.ColorFunc("blue+h:black")
)


var (
	levels = map[string]colorfn{
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
				fmt.Println(color_custom(line))
				continue
			}
		}

		if match := pattern_level.FindString(line); len(match) > 0 {
			if color, ok := levels[strings.TrimSpace(match)]; ok {
				fmt.Println(color(line))
				continue
			}
		}

		fmt.Println(line)
	}
}

func main() {
	var custom = flag.String("custom", "", "custom pattern")
	var custom_color = flag.String("custom-color", "white+h:blue", "custom color (default white on blue")

        flag.Parse()

	if len(*custom) > 1 {
		pattern_custom = regexp.MustCompile(*custom)
	}
	if len(*custom_color) > 1 {
		color_custom = ansi.ColorFunc(*custom_color)
	}

	Colorize(os.Stdin)
}
