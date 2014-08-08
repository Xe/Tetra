package r1459

import (
	"testing"
)

func TestBaseParse(t *testing.T) {
	line := "FOO"

	lineStruct := NewRawLine(line)

	if lineStruct.Verb != "FOO" {
		t.Fatalf("Line verb expected to be FOO, it is %s", lineStruct.Verb)
	}
}

func TestPRIVMSGParse(t *testing.T) {
	line := ":Xena!oper@yolo-swag.com PRIVMSG #niichan :Why hello there"

	lineStruct := NewRawLine(line)

	if lineStruct.Verb != "PRIVMSG" {
		t.Fatalf("Line verb expected to be PRIVMSG, it is %s", lineStruct.Verb)
	}

	if lineStruct.Source != "Xena!oper@yolo-swag.com" {
		t.Fatalf("Line source expected to be PRIVMSG, it is %s", lineStruct.Source)
	}

	if len(lineStruct.Args) != 2 {
		t.Fatalf("Line arg count expected to be 2, it is %s", len(lineStruct.Args))
	}

	if lineStruct.Args[0] != "#niichan" {
		t.Fatalf("Line arg 0 expected to be #niichan, it is %s", lineStruct.Args[0])
	}

	if lineStruct.Args[1] != "Why hello there" {
		t.Fatalf("Line arg 1 expected to be 'Why hello there', it is %s", lineStruct.Args[1])
	}
}

// This test case has previously been known to crash this library.
func TestPreviouslyBreakingLine(t *testing.T) {
	line := ":649AAAABS AWAY"

	lineStruct := NewRawLine(line)

	if lineStruct.Source != "649AAAABS" {
		t.Fatalf("Line source expected to be 649AAAABS, it is %s", lineStruct.Source)
	}

	if lineStruct.Verb != "AWAY" {
		t.Fatalf("Line verb expected to be AWAY, it is %s", lineStruct.Verb)
	}
}
