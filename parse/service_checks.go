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
			return nil
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

				return ServiceInsideCheckProcessMethodsStringValue
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
			return ServiceInsideCheckProcessMethodsStringValue
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
			return ServiceInsideCheckProcessMethodsStringValue
		}
	}

	return nil
}

func ServiceInsideCheckProcessMethodsStringValue(l *lexer) stateFn {
	if l.next() != '"' {
		return l.errorf("check process missing value")
	}

	for {
		switch nextRune := l.next(); {
		case isEndOfLine(nextRune) || isEof(nextRune):
			return l.errorf("check process missing value %s", l.input[l.pos:])
		case nextRune == '"':
			l.emit(itemInsideCheckProcess_ProgramMethodStringValue)
			l.skipWhiteSpaces()
			return ServiceInsideCheckProcessMethods
		}
	}
	return ServiceInsideCheckProcessMethods
}
