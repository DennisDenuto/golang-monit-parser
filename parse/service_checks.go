package lex

import (
	"strings"
	"fmt"
	"errors"
)

func ServiceCheckStart(l *lexer) stateFn {
	l.skipWhiteSpaces()
	l.pos += len("check")
	l.emit(itemCheckStart)

	if x := strings.Index(l.input[l.pos:], "process"); x >= 0 {
		l.skipWhiteSpaces()
		return ServiceCheckProcessStart
	}

	if x := strings.Index(l.input[l.pos:], "file"); x >= 0 {
		l.skipWhiteSpaces()
		return ServiceCheckFileStart
	}

	return nil
}

func ServiceCheckProcessStart(l *lexer) stateFn {
	l.pos += len("process")
	l.emit(itemCheckProcess)

	l.skipWhiteSpaces()

	return ServiceInsideCheckProcess
}

func ServiceCheckFileStart(l *lexer) stateFn {
	l.pos += len("file")
	l.emit(itemCheckFile)

	l.skipWhiteSpaces()

	return ServiceInsideCheckFile
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
			return ServiceInsideCheckProcessPid
		case isEof(nextRune):
			l.emit(itemInsideCheckProcess_Name)
			return nil
		}
	}

	return nil
}
func ServiceInsideCheckFile(l *lexer) stateFn {
	for {
		switch nextRune := l.next(); {
		case isAlphaNumeric(nextRune):
		case isSpace(nextRune):
			l.backup()
			l.emit(itemInsideCheckFile_Name)
			l.skipWhiteSpaces()

			if strings.HasPrefix(l.input[l.pos:], "path") {
				return ServiceInsideCheckPath
			}
			return l.errorf("check file <path> missing")
		case isEndOfLine(nextRune):
			l.backup()
			l.emit(itemInsideCheckFile_Name)
			l.skipWhiteSpaces()
			return nil
		case isEof(nextRune):
			l.emit(itemInsideCheckFile_Name)
			return nil
		}
	}

	return nil
}

func ServiceInsideCheckPath(l *lexer) stateFn {
	l.pos += len("path")
	l.skipWhiteSpaces()

	for {
		switch nextRune := l.next(); {
		case isAlphaNumeric(nextRune):
		case isSpace(nextRune):
			l.backup()
			l.emit(itemInsideCheckFile_Path)
			l.skipWhiteSpaces()
			return ServiceInsideCheckProcessMethods
		case isEndOfLine(nextRune):
			l.backup()
			l.emit(itemInsideCheckFile_Path)
			l.skipWhiteSpaces()
			return ServiceInsideCheckProcessMethods
		case isEof(nextRune):
			l.emit(itemInsideCheckFile_Path)
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
	if strings.HasPrefix(l.input[l.pos:], "start") || strings.HasPrefix(l.input[l.pos:], "stop") {
		localItemInsideCheckProcessProgramMethod := itemInsideCheckProcess_StartProgramMethod

		if strings.HasPrefix(l.input[l.pos:], "stop") {
			localItemInsideCheckProcessProgramMethod = itemInsideCheckProcess_StopProgramMethod
		}

		for {
			switch nextRune := l.next(); {
			case isAlphaNumeric(nextRune):
			case nextRune == '=':
				l.backup()
				for {
					if isSpace(l.current()) {
						l.backup()
					} else {
						break
					}
				}

				l.emit(localItemInsideCheckProcessProgramMethod)
				l.emit(itemInsideCheckProcess_ProgramMethodPath)
				l.acceptRun(" =")
				l.ignore()

				err := emitStringValue(l)
				if err != nil {
					return l.errorf(err.Error())
				}
				return ServiceInsideCheckProcessMethods
			case isEndOfLine(nextRune) || isEof(nextRune):
				return l.errorf("check process start missing '='", l.input[l.start:l.pos])
			}
		}

		return ServiceInsideCheckProcessMethods
	}
	if strings.HasPrefix(l.input[l.pos:], "as") {
		l.pos += len("as")
		l.skipWhiteSpaces()
		l.ignore()
		switch {
		case strings.HasPrefix(l.input[l.pos:], "uid"):
			l.pos += len("uid")
			l.emit(itemInsideCheckProcess_ProgramMethodUid)
			l.skipWhiteSpaces()

			err := emitStringValue(l)
			if err != nil {
				return l.errorf(err.Error())
			}
			return ServiceInsideCheckProcessMethods
		}
	}
	if strings.HasPrefix(l.input[l.pos:], "and") {
		l.pos += len("and")
		l.skipWhiteSpaces()
		l.ignore()
		switch {
		case strings.HasPrefix(l.input[l.pos:], "gid"):
			l.pos += len("gid")
			l.emit(itemInsideCheckProcess_ProgramMethodGid)
			l.skipWhiteSpaces()
			err := emitStringValue(l)
			if err != nil {
				return l.errorf(err.Error())
			}

			return ServiceInsideCheckProcessMethods
		}
	}
	if strings.HasPrefix(l.input[l.pos:], "group") {
		l.pos += len("group")
		l.emit(itemInsideCheckProcess_ProgramMethodGroupName)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}

		return ServiceInsideCheckProcessMethods
	}
	if strings.HasPrefix(l.input[l.pos:], "depends on") {
		l.pos += len("depends on")
		l.emit(itemServiceDependencies)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}

		return ServiceInsideCheckProcessMethods
	}
	if strings.HasPrefix(l.input[l.pos:], "if failed") {
		l.pos += len("if failed")
		l.emit(itemInsideCheckProcess_ConnectionTestingEnterIfConditions)
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessConnectionTesting
	}
	if strings.HasPrefix(l.input[l.pos:], "if total memory") {
		l.pos += len("if total memory")
		l.emit(itemInsideCheckResourceTesting)
		l.skipWhiteSpaces()
		return InsideCheckResourceTesting
	}
	return nil
}

func InsideCheckResourceTesting(l *lexer) stateFn {

	if l.accept("><=") {
		l.accept("><=")
		l.emit(itemInsideCheckResourceTestingOperator)

		l.skipWhiteSpaces()
		l.acceptNumbers()
		l.acceptRun(" ")
		l.acceptRun("%kmgbKMGB")
		l.emit(itemInsideCheckProcess_ProgramMethodUnQuotedStringValue)
		l.skipWhiteSpaces()

		return InsideCheckResourceTesting
	}

	if strings.HasPrefix(l.input[l.pos:], "for ") {
		l.acceptUntilSpace()
		l.emit(itemInsideCheckProcess_ConnectionTesting_Cycle)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		err = emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		l.skipWhiteSpaces()
		return InsideCheckResourceTesting
	}

	return ServiceInsideCheckProcessConnectionTesting
}

func ServiceInsideCheckProcessConnectionTesting(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "unixsocket ") {
		l.acceptUntilSpace()
		l.emit(itemInsideCheckProcess_ConnectionTesting_UnixSocket)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessInsideConnectionTesting
	}

	if strings.HasPrefix(l.input[l.pos:], "host ") {
		l.acceptUntilSpace()
		l.emit(itemInsideCheckProcess_ConnectionTesting_TcpUdpHost)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessConnectionTesting
	}

	if strings.HasPrefix(l.input[l.pos:], "port ") {
		l.acceptUntilSpace()
		l.emit(itemInsideCheckProcess_ConnectionTesting_TcpUdpPort)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessConnectionTesting
	}

	if strings.HasPrefix(l.input[l.pos:], "protocol ") {
		l.acceptUntilSpace()
		l.emit(itemInsideCheckProcess_ConnectionTesting_TcpUdpProtocol)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessConnectionTesting
	}

	if strings.HasPrefix(l.input[l.pos:], "then ") {
		l.acceptUntilSpace()
		l.emit(itemInsideCheckProcess_ConnectionTesting_Action)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessMethods
	}

	return ServiceInsideCheckProcessInsideConnectionTesting
}

func ServiceInsideCheckProcessInsideConnectionTesting(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "with timeout ") {
		l.pos += len("with tineout")
		l.emit(itemInsideCheckProcess_ConnectionTesting_Timeout)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		err = emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		l.skipWhiteSpaces()

		return ServiceInsideCheckProcessInsideConnectionTesting
	}

	if strings.HasPrefix(l.input[l.pos:], "for ") {
		l.acceptUntilSpace()
		l.emit(itemInsideCheckProcess_ConnectionTesting_Cycle)
		l.skipWhiteSpaces()
		err := emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		err = emitStringValue(l)
		if err != nil {
			return l.errorf(err.Error())
		}
		l.skipWhiteSpaces()
		return ServiceInsideCheckProcessInsideConnectionTesting
	}
	l.emit(itemInsideCheckProcess_ConnectionTesting_ExitIfConditions)
	return ServiceInsideCheckProcessConnectionTesting
}

/*
Strings can be either quoted or unquoted. A quoted string is bounded by double quotes and may contain whitespace (and quoted digits are treated as a string). An unquoted string is any whitespace-delimited token, containing characters and/or numbers.
 */
func emitStringValue(l *lexer) error {
	next := l.next()
	if next == '"' {
		for {
			switch nextRune := l.next(); {
			case isEndOfLine(nextRune) || isEof(nextRune):
				return errors.New(fmt.Sprintf("check process missing value %s", l.input[l.pos:]))
			case nextRune == '"':
				l.emit(itemInsideCheckProcess_ProgramMethodQuotedStringValue)
				l.skipWhiteSpaces()
				return nil
			}
		}
	} else if isAlphaNumeric(next) {
		l.acceptUntilSpace()
		l.emit(itemInsideCheckProcess_ProgramMethodUnQuotedStringValue)
		l.skipWhiteSpaces()
		return nil
	}

	return nil
}
