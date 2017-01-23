package lex_test

import (
	. "github.com/DennisDenuto/golang-monit-parser/parse"

	"github.com/DennisDenuto/golang-monit-parser/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parse", func() {
	var parser Parser

	BeforeEach(func() {
		parser = NewMonitParser()
	})

	Context("Monit file with single check process", func() {
		var monitFileContents string
		BeforeEach(func() {
			monitFileContents = `check process abc pidfile /tmp`
		})
		It("should build monit tree with check process", func() {
			_, items := Lex("test", monitFileContents)

			monitFileParsed := parser.Parse(items)
			Expect(monitFileParsed).ToNot(BeNil())
			Expect(monitFileParsed.CheckProcesses).ToNot(BeNil())
			Expect(monitFileParsed.CheckProcesses).To(ConsistOf(
				api.ProcessCheck{
					Name:    "abc",
					Pidfile: "/tmp",
				},
			))
		})

		Context("with service methods", func() {
			BeforeEach(func() {
				monitFileContents += "\n"
				monitFileContents += `  start program = "/usr/local/mmonit/bin/mmonit" as uid "mmonit" and gid "gmmonit"
  stop program = "/usr/local/mmonit/bin/mmonit stop" as uid "mmonit" and gid "gmmonit"`
			})

			It("should build monit tree with check process", func() {
				_, items := Lex("test", monitFileContents)

				monitFileParsed := parser.Parse(items)
				Expect(monitFileParsed).ToNot(BeNil())
				Expect(monitFileParsed.CheckProcesses).ToNot(BeNil())
				Expect(monitFileParsed.CheckProcesses).To(ConsistOf(
					api.ProcessCheck{
						Name:    "abc",
						Pidfile: "/tmp",
						StartProgram: api.CheckProgram{
							Path: "/usr/local/mmonit/bin/mmonit",
							Uid:  "mmonit",
							Gid:  "gmmonit",
						},
						StopProgram: api.CheckProgram{
							Path: "/usr/local/mmonit/bin/mmonit stop",
							Uid:  "mmonit",
							Gid:  "gmmonit",
						},
					},
				))
			})

			Context("with connection testing", func() {
				BeforeEach(func() {
					monitFileContents += "\n"
					monitFileContents += `  if failed unixsocket /path/to/socket.sock
    with timeout 55 seconds for 5 cycles
  then restart`
				})

				It("should build monit tree with connection testing", func() {
					_, items := Lex("test", monitFileContents)

					monitFileParsed := parser.Parse(items)
					Expect(monitFileParsed).ToNot(BeNil())
					Expect(monitFileParsed.CheckProcesses).ToNot(BeNil())
					Expect(monitFileParsed.CheckProcesses).To(ConsistOf(
						api.ProcessCheck{
							Name:    "abc",
							Pidfile: "/tmp",
							StartProgram: api.CheckProgram{
								Path: "/usr/local/mmonit/bin/mmonit",
								Uid:  "mmonit",
								Gid:  "gmmonit",
							},
							StopProgram: api.CheckProgram{
								Path: "/usr/local/mmonit/bin/mmonit stop",
								Uid:  "mmonit",
								Gid:  "gmmonit",
							},
							FailedSocket: api.FailedSocket{
								SocketFile: "/path/to/socket.sock",
								Timeout:    55,
								NumCycles:  5,
								Action:     "restart",
							},
						},
					))
				})

			})
		})
	})

	Context("Monit file with multiple process checks", func() {
		It("should build monit tree with check process", func() {

			monitFileContents := `check process short_process
  pidfile /path/to/short/pid

check process another_process
  with pidfile /path/to/another/pid
  start program = "/path/to/short/start/command"`

			_, items := Lex("test", monitFileContents)

			monitFileParsed := parser.Parse(items)
			Expect(monitFileParsed).ToNot(BeNil())
			Expect(monitFileParsed.CheckProcesses).ToNot(BeNil())
			Expect(monitFileParsed.CheckProcesses).To(ConsistOf(
				api.ProcessCheck{
					Name:    "short_process",
					Pidfile: "/path/to/short/pid",
				},
				api.ProcessCheck{
					Name:         "another_process",
					Pidfile:      "/path/to/another/pid",
					StartProgram: api.CheckProgram{Path: "/path/to/short/start/command"},
				},
			))
		})

	})

})
