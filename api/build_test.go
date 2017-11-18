package api_test

import (
	"fmt"

	"github.com/bronzdoc/pak/api"
	"github.com/bronzdoc/pak/pakfile"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build", func() {
	var pakFile *pakfile.PakFile
	testDir := "/tmp/test_pak"
	artifactName := "pak-test-artifact"

	BeforeEach(func() {
		os.Mkdir(testDir, 0777)

		pakFile = &pakfile.PakFile{}
		pakFile.ArtifactName = artifactName
		pakFile.Path = testDir
		pakFile.Metadata = map[string]string{
			"key": "value",
		}
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
		os.Remove(fmt.Sprintf("%s.tar", artifactName))
	})

	It("should build a pak package", func() {
		api.Build(pakFile)
		_, err := os.Stat(fmt.Sprintf("%s.tar", artifactName))
		Expect(os.IsNotExist(err)).To(Equal(false))
	})
})
