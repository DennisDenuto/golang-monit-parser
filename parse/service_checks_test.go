package lex

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lex/ServiceChecks", func() {
	var act func(string) *lexer

	BeforeEach(func() {
		act = func(input string) *lexer {
			return &lexer{
				name:  "test",
				input: input,
				items: make(chan Item, 10),
			}
		}
	})

	Context("Check Process", func() {
		It("Should scan check process with only name", func() {
			lex := act("check process abc")

			nextLexFn := ServiceCheckStart(lex)
			Expect(lex.pos).To(Equal(6))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckStart, Value: "check"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(14))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckProcess, Value: "process"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(17))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_Name, Value: "abc"})))

			Expect(nextLexFn).To(BeNil())
		})

		It("Should scan check process with process file", func() {
			lex := act("check process abc pidfile /tmp")

			nextLexFn := ServiceCheckStart(lex)
			Expect(lex.pos).To(Equal(6))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckStart, Value: "check"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(14))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckProcess, Value: "process"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(18))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_Name, Value: "abc"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_Pid, Value: "pidfile /tmp"})))
			Expect(lex.pos).To(Equal(30))
		})

		It("Should scan check process with process file using with keyword", func() {
			lex := act(`check process abc
  with pidfile /tmp`)

			nextLexFn := ServiceCheckStart(lex)
			Expect(lex.pos).To(Equal(6))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckStart, Value: "check"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(14))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckProcess, Value: "process"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(20))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_Name, Value: "abc"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_Pid, Value: "with pidfile /tmp"})))
			Expect(lex.pos).To(Equal(37))
		})

		It("Should scan check process with process regex", func() {
			lex := act("check process abc matching foobar.*")

			nextLexFn := ServiceCheckStart(lex)
			Expect(lex.pos).To(Equal(6))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckStart, Value: "check"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(14))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckProcess, Value: "process"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(18))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_Name, Value: "abc"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_Pid, Value: "matching foobar.*"})))
			Expect(lex.pos).To(Equal(35))
		})

		Context("With service methods", func() {
			It("should scan check process with service methods", func() {
				lex := act(`check process abc matching foobar.*
  start program = "/usr/local/mmonit/bin/mmonit" as uid "mmonit" and gid "mmonit"
  stop program = "/usr/local/mmonit/bin/mmonit stop" as uid "stop_mmonit" and gid "stop_mmonit"
  group group_name
  depends on file_check
  if total memory > 2048 Mb for 3 cycles then alert`)

				nextLexFn := ServiceCheckStart(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_StartProgramMethod, Value: `start program`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodPath, Value: ""})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"/usr/local/mmonit/bin/mmonit"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUid, Value: "uid"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"mmonit"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodGid, Value: "gid"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"mmonit"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_StopProgramMethod, Value: `stop program`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodPath, Value: ""})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"/usr/local/mmonit/bin/mmonit stop"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUid, Value: "uid"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"stop_mmonit"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodGid, Value: "gid"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"stop_mmonit"`})))

				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodGroupName, Value: "group"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: "group_name"})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemServiceDependencies, Value: "depends on"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: "file_check"})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckResourceTesting, Value: "if total memory"})))


				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckResourceTestingOperator, Value: ">"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: "2048 Mb"})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)

				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_Cycle, Value: "for"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `3`})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `cycles`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_Action, Value: "then"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `alert`})))
			})
		})

		Context("With connection testing", func() {
			It("should scan check process with socket test", func() {
				lex := act(`check process abc matching foobar.*
  if failed unixsocket /path/to/socket.sock
    with timeout 5 seconds for 5 cycles
  then restart`)

				/*
FOR <X> CYCLES ...
or:
 <X> [TIMES WITHIN] <Y> CYCLES ...

 IF FAILED
    <UNIXSOCKET path>
    [TYPE <TCP|UDP>]
    [PROTOCOL protocol | <SEND|EXPECT> "string",...]
    [TIMEOUT number SECONDS]
    [RETRY number]
 THEN action
 */

				nextLexFn := ServiceCheckStart(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTestingEnterIfConditions, Value: "if failed"})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_UnixSocket, Value: "unixsocket"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `/path/to/socket.sock`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_Timeout, Value: "with timeout"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `5`})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `seconds`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_Cycle, Value: "for"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `5`})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `cycles`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_ExitIfConditions, Value: ""})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_Action, Value: "then"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `restart`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).To(BeNil())
			})

			It("should scan check process with host test", func() {
				lex := act(`check process abc matching foobar.*
  if failed host 1.2.3.4 port 9876 protocol http with timeout 20 seconds for 10 cycles
  then stop`)

				nextLexFn := ServiceCheckStart(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive())

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTestingEnterIfConditions, Value: "if failed"})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_TcpUdpHost, Value: "host"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `1.2.3.4`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_TcpUdpPort, Value: "port"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `9876`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_TcpUdpProtocol, Value: "protocol"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `http`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)

				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_Timeout, Value: "with timeout"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `20`})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `seconds`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).ToNot(BeNil())

				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_Cycle, Value: "for"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `10`})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `cycles`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_ExitIfConditions, Value: ""})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ConnectionTesting_Action, Value: "then"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUnQuotedStringValue, Value: `stop`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).To(BeNil())
			})
		})
	})

	Context("Check File", func() {
		It("Should scan check file", func() {
			lex := act("check file unique-name path /tmp/test")

			nextLexFn := ServiceCheckStart(lex)
			Expect(lex.pos).To(Equal(6))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckStart, Value: "check"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(11))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemCheckFile, Value: "file"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.pos).To(Equal(23))
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckFile_Name, Value: "unique-name"})))

			Expect(nextLexFn).ToNot(BeNil())
			nextLexFn = nextLexFn(lex)
			Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckFile_Path, Value: "/tmp/test"})))
			Expect(nextLexFn).To(BeNil())
		})

		Context("With service methods", func() {
			It("should scan check file with service methods", func() {
				lex := act(`check file unique-name path /tmp/test
  start program = "/usr/local/mmonit/bin/mmonit" as uid "mmonit" and gid "mmonit"
  stop program = "/usr/local/mmonit/bin/mmonit stop" as uid "stop_mmonit" and gid "stop_mmonit"`)

				nextLexFn := ServiceCheckStart(lex)

				var nextItem Item
				Expect(lex.items).To(Receive(&nextItem))
				Expect(nextItem.Type).To(Equal(itemCheckStart))

				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).ToNot(BeNil())
				Expect(lex.items).To(Receive(&nextItem))
				Expect(nextItem.Type).To(Equal(itemCheckFile))

				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).ToNot(BeNil())
				Expect(lex.items).To(Receive(&nextItem))
				Expect(nextItem.Type).To(Equal(itemInsideCheckFile_Name))

				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).ToNot(BeNil())
				Expect(lex.items).To(Receive(&nextItem))
				Expect(nextItem.Type).To(Equal(itemInsideCheckFile_Path))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_StartProgramMethod, Value: `start program`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodPath, Value: ""})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"/usr/local/mmonit/bin/mmonit"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUid, Value: "uid"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"mmonit"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodGid, Value: "gid"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"mmonit"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_StopProgramMethod, Value: `stop program`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodPath, Value: ""})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"/usr/local/mmonit/bin/mmonit stop"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodUid, Value: "uid"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"stop_mmonit"`})))

				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodGid, Value: "gid"})))
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_ProgramMethodQuotedStringValue, Value: `"stop_mmonit"`})))

				Expect(nextLexFn).To(BeNil())
			})
		})
	})
})
