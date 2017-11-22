package api_test

import (
	"archive/tar"
	"fmt"
	"io/ioutil"

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
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
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

	It("should create a pak.metadata inside the package", func() {
		metadataFileName := "pak.metadata"
		api.Build(pakFile)

		pkg, _ := os.Open(fmt.Sprintf("%s.tar", artifactName))

		tr := tar.NewReader(pkg)

		// Search for the metadata file inside the package
		metadataFileExist := func() bool {
			for {
				header, _ := tr.Next()
				if header == nil {
					return false
				}

				if header.Name == metadataFileName {
					return true
				}
			}
		}()

		content, _ := ioutil.ReadAll(tr)

		expectedContent := `{
  "build": {
    "key1": "value1",
    "key2": "value2",
    "key3": "value3"
  }
}`
		Expect(metadataFileExist).To(BeTrue())
		//Expect(string(content)).To(Equal("key1=value1\nkey2=value2\nkey3=value3\n"))
		Expect(string(content)).To(Equal(expectedContent))
	})
})
