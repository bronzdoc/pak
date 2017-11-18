package api_test

import (
	"fmt"
	"os"

	. "github.com/bronzdoc/pak/api"
	"github.com/jhoonb/archivex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inspect", func() {
	testDir := "/tmp/test_pak"
	artifactName := "pak-test-artifact.tar"

	BeforeEach(func() {
		os.Mkdir(testDir, 0777)
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	It("should inspect an artifact", func() {
		artifact := new(archivex.TarFile)
		artifact.Create(fmt.Sprintf("%s/%s", testDir, artifactName))

		keyValPair := "key=value"
		artifact.Add("pak.metadata", []byte(keyValPair))

		artifact.Close()

		content, err := Inspect(fmt.Sprintf("%s/%s", testDir, artifactName))
		Expect(err).To(BeNil())
		Expect(content).To(Equal("key=value"))
	})

	Context("when no metadata found ", func() {
		It("should return the correct error message", func() {
			artifact := new(archivex.TarFile)
			artifact.Create(fmt.Sprintf("%s/%s", testDir, artifactName))
			artifact.Close()

			content, err := Inspect(fmt.Sprintf("%s/%s", testDir, artifactName))
			Expect(err.Error()).To(Equal("no metadata found"))
			Expect(content).To(Equal(""))
		})
	})
})
