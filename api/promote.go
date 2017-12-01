package api

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/bronzdoc/pak/util"
	"github.com/jhoonb/archivex"
	"github.com/pkg/errors"
)

func Promote(packageName string, options map[string]interface{}) error {
	inspectOptions := map[string]interface{}{}

	// get metadata
	jsonMetadata, err := Inspect(packageName, inspectOptions)
	if err != nil {
		return errors.Wrap(err, "failed to inspect package metadata")
	}

	var metadata map[string]interface{}
	json.Unmarshal([]byte(jsonMetadata), &metadata)

	// get metadata of the label we are promoting
	var labelMetadata map[string]interface{}
	var label string
	if l, ok := options["label"]; ok {
		labelMetadata = metadata[l.(string)].(map[string]interface{})
		label = l.(string)
	}

	// Resolve env variables of the metadata
	for key, value := range labelMetadata {
		var newValue string
		var err error

		switch value.(type) {
		case string:
			newValue, err = util.ResolveEnvVar(value.(string))
			labelMetadata[key] = newValue
		case map[string]interface{}:
			resolvedMetadata := make(map[string]interface{})
			for k, v := range value.(map[string]interface{}) {
				newValue, err = util.ResolveEnvVar(v.(string))
				resolvedMetadata[k] = newValue
			}

			labelMetadata[key] = resolvedMetadata
		}

		if err != nil {
			return errors.Wrap(err, "failed to resolve env var")
		}
	}

	metadata[label] = labelMetadata

	// Create promote artifact
	var newArtifactname string
	if value, ok := labelMetadata["name"]; ok {
		newArtifactname = value.(string)
	}

	newArtifact := new(archivex.TarFile)
	newArtifact.Create(fmt.Sprintf("%s.tar", newArtifactname))

	// create new pak.metadata with the old data and the new one
	jsonString, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		panic(err)
	}

	// Store metadata in the package
	newArtifact.Add("pak.metadata", jsonString)

	// Add packaged files form the old artifact to the new one
	pkg, err := os.Open(packageName)
	if err != nil {
		return errors.Wrapf(err, "failed to open package %s", packageName)
	}

	tr := tar.NewReader(pkg)

	// Search for the metadata file inside the package
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return errors.Wrap(err, "could not get tar header")
		}

		content, err := ioutil.ReadAll(tr)
		if err != nil {
			return errors.Wrapf(err, "could not read package %s", packageName)
		}

		// Don't add metadata file since we will be adding a new one with new metadata
		metadataFileName := "pak.metadata"
		if header.Name != metadataFileName {
			newArtifact.Add(header.Name, content)
		}
	}

	newArtifact.Close()

	return nil
}

// Recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	//if !fi.IsDir() {
	//	return &CustomError{"Source is not a directory"}
	//}

	// ensure dest dir does not already exist

	_, err = os.Open(dest)
	//if !os.IsNotExist(err) {
	//	return &CustomError{"Destination already exists"}
	//}

	// create dest dir

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {

		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return
}

// Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return
}
