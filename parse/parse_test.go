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
		It("should build monit tree with check process", func() {
			monitFileContents := `check process abc pidfile /tmp`
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
			It("should build monit tree with check process", func() {
				monitFileContents := `check process abc pidfile /tmp
  start program = "/usr/local/mmonit/bin/mmonit" as uid "mmonit" and gid "gmmonit"
  stop program = "/usr/local/mmonit/bin/mmonit stop" as uid "mmonit" and gid "gmmonit"`
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
		})
	})
})
