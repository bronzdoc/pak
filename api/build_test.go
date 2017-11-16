package api_test

import (
	"github.com/bronzdoc/pak/api"
	"github.com/bronzdoc/pak/pakfile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Build", func() {
	var pakFile *pakfile.PakFile

	BeforeEach(func() {
		os.Mkdir("/tmp/test_pac", 0777)

		pakFile = &pakfile.PakFile{}
		pakFile.ArtifactName = "pak-test-artifact"
		pakFile.Path = "/tmp/test_pac"
		pakFile.Metadata = map[string]string{
			"key": "value",
		}
	})

	AfterEach(func() {
		os.Remove("./tmp/test_pac")
		os.Remove("pak-test-artifact.tar.gz")
	})

	It("should build a pak package", func() {
		api.Build(pakFile)
		_, err := os.Stat("pak-test-artifact.tar.gz")
		Expect(os.IsNotExist(err)).To(Equal(false))
	})
})
