package lex

import (
	"github.com/DennisDenuto/golang-monit-parser/api"
	"strings"
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
			monitFileParsed.CheckProcesses[0].Name = item.Value
		case itemInsideCheckProcess_Pid:
			monitFileParsed.CheckProcesses[0].Pidfile = item.Value[8:]
		case itemInsideCheckProcess_StartProgramMethod:
			<-items
			pathValue := <-items
			<-items
			uid := <-items
			<-items
			gid := <-items
			monitFileParsed.CheckProcesses[0].StartProgram = api.CheckProgram{
				Path: strings.Replace(pathValue.Value, `"`, "", -1),
				Uid:  strings.Replace(uid.Value, `"`, "", -1),
				Gid:  strings.Replace(gid.Value, `"`, "", -1),
			}
		case itemInsideCheckProcess_StopProgramMethod:
			<-items
			pathValue := <-items
			<-items
			uid := <-items
			<-items
			gid := <-items
			monitFileParsed.CheckProcesses[0].StopProgram = api.CheckProgram{
				Path: strings.Replace(pathValue.Value, `"`, "", -1),
				Uid:  strings.Replace(uid.Value, `"`, "", -1),
				Gid:  strings.Replace(gid.Value, `"`, "", -1),
			}

		}
	}

	return monitFileParsed
}

type MonitFileParsed struct {
	CheckProcesses []api.ProcessCheck
}
