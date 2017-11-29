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

		jsonString := `{"key":"value"}`
		artifact.Add("pak.metadata", []byte(jsonString))

		artifact.Close()

		var options map[string]interface{}
		content, err := Inspect(
			fmt.Sprintf("%s/%s", testDir, artifactName),
			options,
		)

		Expect(err).To(BeNil())
		Expect(content).To(Equal(`{"key":"value"}`))
	})

	Context("when is_key_value options is passed", func() {
		It("should return content as key value pair", func() {
			artifact := new(archivex.TarFile)
			artifact.Create(fmt.Sprintf("%s/%s", testDir, artifactName))

			jsonString := `{"label":{"key":"value"}}`
			artifact.Add("pak.metadata", []byte(jsonString))

			artifact.Close()

			options := map[string]interface{}{
				"is_key_value": true,
			}

			content, err := Inspect(
				fmt.Sprintf("%s/%s", testDir, artifactName),
				options,
			)

			Expect(err).To(BeNil())
			Expect(content).To(Equal("#label\n  key=\"value\"\n"))
		})
	})

	Context("when a specific metadata label is given", func() {
		It("should inspect only inspect only a subset of the metadata", func() {
			artifact := new(archivex.TarFile)
			artifact.Create(fmt.Sprintf("%s/%s", testDir, artifactName))

			jsonString := `{"label":{"key":"value"}, "test_label":{"test_key":"test_value"}}`
			artifact.Add("pak.metadata", []byte(jsonString))

			artifact.Close()

			options := map[string]interface{}{
				"label": "test_label",
			}

			content, err := Inspect(
				fmt.Sprintf("%s/%s", testDir, artifactName),
				options,
			)

			Expect(err).To(BeNil())
			Expect(content).To(Equal("{\n  \"test_key\": \"test_value\"\n}"))
		})
	})

	Context("when no metadata found", func() {
		It("should return the correct error message", func() {
			artifact := new(archivex.TarFile)
			artifact.Create(fmt.Sprintf("%s/%s", testDir, artifactName))
			artifact.Close()

			var options map[string]interface{}
			content, err := Inspect(
				fmt.Sprintf("%s/%s", testDir, artifactName),
				options,
			)

			Expect(err.Error()).To(Equal("no metadata found"))
			Expect(content).To(Equal(""))
		})
	})
})
