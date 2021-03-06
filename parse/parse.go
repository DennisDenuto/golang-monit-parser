package lex

import (
	"github.com/DennisDenuto/golang-monit-parser/api"
	"strings"
	"strconv"
)

type Parser struct{}

func NewMonitParser() Parser {
	return Parser{}
}

func (Parser) Parse(items chan Item) MonitFileParsed {
	monitFileParsed := MonitFileParsed{}

	for item := range items {
		switch item.Type {
		case itemCheckProcess:
			monitFileParsed.CheckProcesses = append(monitFileParsed.CheckProcesses, api.ProcessCheck{})
		case itemInsideCheckProcess_Name:
			check := monitFileParsed.CheckProcesses.GetLast()
			check.Name = item.Value
		case itemInsideCheckProcess_Pid:
			check := monitFileParsed.CheckProcesses.GetLast()
			check.Pidfile = removeNoiseKeyword(item.Value)[len("pidfile "):]
		case itemInsideCheckProcess_StartProgramMethod:
			<-items
			pathValue := <-items
			<-items
			uid := <-items
			<-items
			gid := <-items
			check := monitFileParsed.CheckProcesses.GetLast()
			check.StartProgram = api.CheckProgram{
				Path: stripQuotes(pathValue.Value),
				Uid:  stripQuotes(uid.Value),
				Gid:  stripQuotes(gid.Value),
			}
		case itemInsideCheckProcess_StopProgramMethod:
			<-items
			pathValue := <-items
			<-items
			uid := <-items
			<-items
			gid := <-items
			check := monitFileParsed.CheckProcesses.GetLast()
			check.StopProgram = api.CheckProgram{
				Path: stripQuotes(pathValue.Value),
				Uid:  stripQuotes(uid.Value),
				Gid:  stripQuotes(gid.Value),
			}
		case itemInsideCheckProcess_ConnectionTestingEnterIfConditions:
			nextItem := <-items
			switch nextItem.Type {
			case itemInsideCheckProcess_ConnectionTesting_UnixSocket:
				socketFilePath := <-items
				<-items
				timeout := <-items
				<-items
				timeoutValue, _ := strconv.Atoi(timeout.Value)

				<-items
				cycle := <-items
				<-items
				cycleValue, _ := strconv.Atoi(cycle.Value)

				<-items
				<-items
				action := <-items

				monitFileParsed.CheckProcesses.GetLast().FailedSocket.SocketFile = socketFilePath.Value
				monitFileParsed.CheckProcesses.GetLast().FailedSocket.Timeout = timeoutValue
				monitFileParsed.CheckProcesses.GetLast().FailedSocket.NumCycles = cycleValue
				monitFileParsed.CheckProcesses.GetLast().FailedSocket.Action = action.Value
			}

		}

	}

	return monitFileParsed
}

/*
noise keywords like 'if', 'and', 'with(in)', 'has', 'us(ing|e)', 'on(ly)', 'then', 'for', 'of'
 */
func removeNoiseKeyword(val string) string {
	return strings.TrimPrefix(val, "with ")
}

func stripQuotes(val string) string {
	return strings.Replace(val, `"`, "", -1)
}

type ProcessChecks []api.ProcessCheck

func (pc ProcessChecks) GetLast() *api.ProcessCheck {
	return &pc[len(pc) - 1]
}

type MonitFileParsed struct {
	CheckProcesses ProcessChecks
}
