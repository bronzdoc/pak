package api_test

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bronzdoc/pak/api"
	"github.com/bronzdoc/pak/pakfile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Promote", func() {
	var pakFile *pakfile.PakFile
	testDir := "/tmp/test_pak"
	artifactName := "pak-test-artifact"
	promoteArtifactName := "promote-artifact"

	BeforeEach(func() {
		os.Setenv("RC_KEY_1", "RC_VALUE_1")
		os.Setenv("RC_KEY_2", "RC_VALUE_2")

		os.Mkdir(testDir, 0777)

		pakFile = &pakfile.PakFile{}
		pakFile.ArtifactName = artifactName
		pakFile.Path = testDir
		pakFile.Metadata = map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		pakFile.Promote = map[string]map[string]interface{}{
			"rc": map[string]interface{}{
				"name": promoteArtifactName,
				"metadata": map[string]interface{}{
					"rc-key1": "${RC_KEY_1}",
					"rc-key2": "${RC_KEY_2}",
					"rc-key3": "RC_KEY_3",
				},
			},
		}

		api.Build(pakFile)
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
		os.Remove(fmt.Sprintf("%s.tar", artifactName))
		os.Remove(fmt.Sprintf("%s.tar", promoteArtifactName))
	})

	It("should create a new artifact", func() {
		err := api.Promote(
			fmt.Sprintf("%s.tar", artifactName),
			map[string]interface{}{
				"label": "rc",
			})

		Expect(err).To(BeNil())

		_, err = os.Stat(fmt.Sprintf("%s.tar", promoteArtifactName))
		Expect(os.IsNotExist(err)).To(Equal(false))
	})

	It("should create a pak.metadata inside the artifact", func() {
		err := api.Promote(
			fmt.Sprintf("%s.tar", artifactName),
			map[string]interface{}{
				"label": "rc",
			})

		Expect(err).To(BeNil())
		metadataFileName := "pak.metadata"

		pkg, _ := os.Open(fmt.Sprintf("%s.tar", promoteArtifactName))

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
    "metadata": {
      "key1": "value1",
      "key2": "value2",
      "key3": "value3"
    },
    "name": "pak-test-artifact"
  },
  "rc": {
    "metadata": {
      "rc-key1": "RC_VALUE_1",
      "rc-key2": "RC_VALUE_2",
      "rc-key3": "RC_KEY_3"
    },
    "name": "promote-artifact"
  }
}`
		Expect(metadataFileExist).To(BeTrue())
		Expect(string(content)).To(Equal(expectedContent))
	})
})
