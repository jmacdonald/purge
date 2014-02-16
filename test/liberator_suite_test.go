package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLiberator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Liberator Suite")
}
