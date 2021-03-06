package api

type ProcessCheck struct {
	Name           string
	Pidfile        string
	StartProgram   CheckProgram
	StopProgram    CheckProgram
	FailedSocket   FailedSocket
	FailedHost     FailedHost
	TotalMemChecks []MemUsage
	Group          string
	DependsOn      string
}


type CheckProgram struct {
	Path string
	Uid string
	Gid string
}

type FileCheck struct {
	Name           string
	Path           string
	IfChanged      string
	FailedSocket   FailedSocket
	FailedHost     FailedHost
	TotalMemChecks []MemUsage
	Group          string
	DependsOn      string
}

type FailedSocket struct {
	SocketFile string
	Timeout    int
	NumCycles  int
	Action     string
}

type FailedHost struct {
	Host      string
	Port      string
	Protocol  string
	Timeout   int
	NumCycles int
	Action    string
}

type MemUsage struct {
	MemLimit  int
	NumCycles int
	Action    string
}