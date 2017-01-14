package lex

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Lex/ServiceChecks", func() {
	var act func(string) *lexer

	BeforeEach(func() {
		act = func(input string) *lexer {
			return &lexer{
				name:  "test",
				input: input,
				items: make(chan Item, 2),
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
  stop program = "/usr/local/mmonit/bin/mmonit stop" as uid "mmonit" and gid "mmonit"`)

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
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_StartMethod, Value: `start program = "/usr/local/mmonit/bin/mmonit" as uid "mmonit" and gid "mmonit"`})))

				Expect(nextLexFn).ToNot(BeNil())
				nextLexFn = nextLexFn(lex)
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_StopMethod, Value: `stop program = "/usr/local/mmonit/bin/mmonit stop" as uid "mmonit" and gid "mmonit"`})))

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
  stop program = "/usr/local/mmonit/bin/mmonit stop" as uid "mmonit" and gid "mmonit"`)

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

				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).ToNot(BeNil())
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_StartMethod, Value: `start program = "/usr/local/mmonit/bin/mmonit" as uid "mmonit" and gid "mmonit"`})))

				nextLexFn = nextLexFn(lex)
				Expect(nextLexFn).ToNot(BeNil())
				Expect(lex.items).To(Receive(Equal(Item{Type: itemInsideCheckProcess_StopMethod, Value: `stop program = "/usr/local/mmonit/bin/mmonit stop" as uid "mmonit" and gid "mmonit"`})))
			})
		})
	})

})
