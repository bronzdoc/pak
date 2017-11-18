package api

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func Inspect(packageName string) (string, error) {
	const metadataFileName string = "pak.metadata"

	pkg, err := os.Open(packageName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open package %s", packageName)
	}

	tr := tar.NewReader(pkg)

	// Search for the metadata file inside the package
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return "", errors.Wrap(err, "could not get tar header")
		}

		if header.Name == metadataFileName {
			break
		}
	}

	content, err := ioutil.ReadAll(tr)
	if err != nil {
		return "", errors.Wrapf(err, "could not read package %s", packageName)
	}

	if len(content) <= 0 {
		return "", fmt.Errorf("no metadata found")
	}

	return string(content), nil
}
