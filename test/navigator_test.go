package liberator_test

import (
	"github.com/jmacdonald/liberator/filesystem/directory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Navigator", func() {
	var (
		navigator *directory.Navigator
		path      string
	)

	BeforeEach(func() {
		navigator = new(directory.Navigator)
	})

	Describe("ChangeDirectory", func() {
		BeforeEach(func() {
			path, _ = os.Getwd()
			navigator.ChangeDirectory(path)
		})

		It("updates CurrentPath with its path argument", func() {
			Expect(navigator.CurrentPath).To(Equal(path))
		})

		It("updates Entries using path argument", func() {
			Expect(navigator.Entries).To(Equal(directory.Entries(path)))
		})
	})
})
