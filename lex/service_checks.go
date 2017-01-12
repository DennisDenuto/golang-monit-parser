package lex

import (
	"strings"
)

func ServiceCheckStart(l *lexer) stateFn {
	l.pos += len("check")
	l.emit(itemCheckStart)

	if x := strings.Index(l.input[l.pos:], "process"); x >= 0 {
		l.skipWhiteSpaces()
		return ServiceCheckProcessStart
	}
	return nil
}

func ServiceCheckProcessStart(l *lexer) stateFn {
	l.pos += len("process")
	l.emit(itemCheckProcess)

	l.skipWhiteSpaces()

	return ServiceInsideCheckProcess
}

func ServiceInsideCheckProcess(l *lexer) stateFn {
	for {
		switch nextRune := l.next(); {
		case isAlphaNumeric(nextRune):
		case isSpace(nextRune):
			l.backup()
			l.emit(itemInsideCheckProcess_Name)
			l.skipWhiteSpaces()
			return ServiceInsideCheckProcessPid
		case isEndOfLine(nextRune):
			l.backup()
			l.emit(itemInsideCheckProcess_Name)
			l.skipWhiteSpaces()
			return nil
		case isEof(nextRune):
			l.emit(itemInsideCheckProcess_Name)
			return nil
		}
	}

	return nil
}
func ServiceInsideCheckProcessPid(l *lexer) stateFn {
	for {
		switch nextRune := l.next(); {
		case isEndOfLine(nextRune):
			l.backup()
			l.emit(itemInsideCheckProcess_Pid)
			l.next()
			l.ignore()

			if isSpace(l.next()) {
				l.skipWhiteSpaces()
				return ServiceInsideCheckProcessMethods
			}
			l.backup()
			return ServiceCheckStart
		case isEof(nextRune):
			l.emit(itemInsideCheckProcess_Pid)
			return nil
		}
	}
	return nil
}

func ServiceInsideCheckProcessMethods(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "start") {
		l.acceptUntilEndOfLine()
		l.emit(itemInsideCheckProcess_StartMethod)
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessMethods
	}
	if strings.HasPrefix(l.input[l.pos:], "stop") {
		l.acceptUntilEndOfLine()
		l.emit(itemInsideCheckProcess_StopMethod)
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessMethods
	}
	return nil
}
