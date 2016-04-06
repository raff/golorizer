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

type pattern_color struct {
	pattern *regexp.Regexp
	color   colorfn
}

//
// Concatenate pattern values and build
// a single regexp to find any of them (pattern1 or pattern2 or ...)
//
func makePatternLevels(patternMap map[string]colorfn) *regexp.Regexp {
	keys := make([]string, 0, len(patternMap))
	for key, _ := range patternMap {
		keys = append(keys, key)
	}

	return regexp.MustCompile("(^|[ \t]+)(" + strings.Join(keys, "|") + ") ")
}

const (
	DEFAULT_CUSTOM_COLOR = "white+h:blue"
)

var (
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
	pattern_custom = make([]pattern_color, 0, 0)
)

//
// this is the main colorizer method, reads from reader and apply colors for matched patterns
//
func Colorize(reader io.Reader, withLevels bool) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		var color colorfn

		for _, pc := range pattern_custom {
			if match := pc.pattern.FindString(line); len(match) > 0 {
				color = pc.color
				break
			}
		}

		if color == nil && withLevels {
			if match := pattern_level.FindString(line); len(match) > 0 {
				color = levels[strings.TrimSpace(match)]
			}
		}

		if color != nil {
			fmt.Println(color(line))
		} else {
			fmt.Println(line)
		}
	}
}

//////////////////////////////////////////////////////////////////////////////
//
// this is used to keep track of the current color to use for custom patterns
//
// for every -color={color} option it update the current color
//
type CurrentColor struct {
	color     string
	colorfunc colorfn
}

func (current *CurrentColor) Set(value string) error {
	if current.color != value {
		if !strings.Contains(value, ":") { // foreground:background
			if value == "black" {
				value += ":white"
			} else {
				value += ":black"
			}
		}

		current.color = value
		current.colorfunc = ansi.ColorFunc(current.color)
	}

	return nil
}

func (current *CurrentColor) String() string {
	return current.color
}

func DefaultColor() CurrentColor {
	color := CurrentColor{}
	color.Set(DEFAULT_CUSTOM_COLOR)
	return color
}

//////////////////////////////////////////////////////////////////////////////
//
// this is used to add custom patterns
//
// for every -pattern={pattern} adds a tuple (compiled-pattern, current-color)
// to pattern_custom
//
// Note that the String() function doesn't return any value (no default)
//
type Custom struct {
	color *CurrentColor
}

func (custom *Custom) Set(value string) error {
	if pattern, err := regexp.Compile(value); err == nil {
		pattern_custom = append(pattern_custom, pattern_color{pattern, custom.color.colorfunc})
		return nil
	} else {
		return err
	}
}

func (custom *Custom) String() (_ string) {
	return
}

//////////////////////////////////////////////////////////////////////////////
//
// program entrypoint
//
func main() {
	color := DefaultColor()
	custom := Custom{&color}

	withLevels := flag.Bool("levels", true, "enable/disable coloring of log levels")
	flag.Var(&color, "color", "custom color")
	flag.Var(&custom, "custom", "custom pattern")
	flag.Parse()

	if len(flag.Args()) == 0 {
		Colorize(os.Stdin, *withLevels)
	} else {
		for _, fname := range flag.Args() {
			if f, err := os.Open(fname); err != nil {
				fmt.Println(err)
				continue
			} else {
				defer f.Close()
				Colorize(f, *withLevels)
			}
		}
	}
}
