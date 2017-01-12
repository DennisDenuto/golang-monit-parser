package lex_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lex Suite")
}
