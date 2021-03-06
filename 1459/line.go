// Package r1459 implements a base structure to scrape out and utilize an RFC 1459
// frame in high level Go code.
package r1459

import (
	"fmt"
	"strings"
)

// RawLine represents an IRC line.
type RawLine struct {
	Source string            `json:"source"`
	Verb   string            `json:"verb"`
	Args   []string          `json:"args"`
	Tags   map[string]string `json:"tags"`
	Raw    string            `json:"-"` // Deprecated
}

// NewRawLine creates a new line and split out an RFC 1459 frame to a RawLine. This will
// not return an error if it fails.
func NewRawLine(input string) (line *RawLine) {
	line = &RawLine{
		Raw: input,
	}

	split := strings.Split(input, " ")

	if split[0][0] == ':' {
		line.Source = split[0][1:]
		line.Verb = split[1]
		split = split[2:]
	} else {
		line.Source = ""
		line.Verb = split[0]
		split = split[1:]
	}

	argstring := strings.Join(split, " ")
	extparam := strings.Split(argstring, " :")

	if len(extparam) > 1 {
		ext := strings.Join(extparam[1:], " :")
		args := strings.Split(extparam[0], " ")

		line.Args = append(args, ext)
	} else {
		line.Args = split
	}

	if len(line.Args) == 0 {
		line.Args = []string{""}
	} else if line.Args[0][0] == ':' {
		line.Args[0] = strings.TrimPrefix(line.Args[0], ":")
	}

	return
}

// String returns the serialized form of a RawLine as an RFC 1459 frame.
func (r *RawLine) String() (res string) {
	if r.Source != "" {
		res = res + fmt.Sprintf(":%s ", r.Source)
	}

	res = res + fmt.Sprintf("%s", r.Verb)

	for i, arg := range r.Args {
		res = res + " "

		if i == len(r.Args)-1 { // Make the last part of the line an extparam
			res = res + ":"
		}

		res = res + arg
	}

	return
}
