package lex

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lex", func() {

	var act func(string) (chan Item)

	BeforeEach(func() {
		act = func(input string) (chan Item) {
			_, tokenChannel := Lex("test", input)
			return tokenChannel
		}

	})
	Context("Check Process", func() {
		It("Should emit check process token", func() {

			items := act("check process abc")

			Eventually(items).Should(Receive(Equal(Item{Type: itemCheckProcess})))
		})
	})

})
