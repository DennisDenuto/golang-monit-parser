package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int
type Pos int

// Item represents a token returned from the scanner.
// Item represents a token or text string returned from the scanner.
type Item struct {
	Type  itemType // The type of this Item.
	Value string   // The value of this Item.
}

const (
	eof        = -1
	spaceChars = " \t\r\n" // These are the space characters defined by Go itself.
)

const (
	itemError itemType = iota // error occurred;
	itemEOF
	itemStringValue

	itemCheckStart

	itemCheckProcess
	itemCheckFile

	itemInsideCheckProcess_Name
	itemInsideCheckProcess_Pid
	itemInsideCheckProcess_ProgramMethodQuotedStringValue
	itemInsideCheckProcess_ProgramMethodUnQuotedStringValue

	itemInsideCheckProcess_StartProgramMethod
	itemInsideCheckProcess_ProgramMethodUid
	itemInsideCheckProcess_ProgramMethodGid
	itemInsideCheckProcess_ProgramMethodGroupName

	itemInsideCheckProcess_StopProgramMethod
	itemInsideCheckProcess_ProgramMethodPath

	itemInsideCheckFile_Name
	itemInsideCheckFile_Path
)

func (i Item) String() string {
	switch i.Type {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.Value
	}
	if len(i.Value) > 10 {
		return fmt.Sprintf("%.10q...", i.Value)
	}
	return fmt.Sprintf("%q", i.Value)
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

type lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this Item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan Item // channel of scanned items.
}

func Lex(name, input string) (*lexer, chan Item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan Item),
	}
	go l.run() // Concurrently run state machine.
	return l, l.items
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
	for state := ServiceCheckStart; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// next returns the next rune in the input.
func (l *lexer) next() (rune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += l.width
	return rune
}

func (l *lexer) current() (rune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos-1 : l.pos])
	return rune
}

// emit passes an Item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- Item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// leftTrimLength returns the length of the spaces at the beginning of the string.
func leftTrimLength(s string) int {
	return len(s) - len(strings.TrimLeft(s, spaceChars))
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) acceptUntilEndOfLine() {
	for {
		next := l.next()
		if isEndOfLine(next) || isEof(next) {
			break
		}
	}
	l.backup()
}

func (l *lexer) acceptUntilSpace() {
	for {
		next := l.next()
		if isSpace(next) || isEof(next) {
			break
		}
	}
	l.backup()
}

func (l *lexer) skipWhiteSpaces() {
	l.pos += leftTrimLength(l.input[l.pos:])
	l.ignore()
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

func isEof(r rune) bool {
	return r == eof
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
