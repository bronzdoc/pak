package util

import (
	"os"
	"regexp"

	"github.com/pkg/errors"
)

func ResolveEnvVar(userEnvVar string) (string, error) {
	IsEnvVar, err := match(userEnvVar, `^\${.+}`)
	if err != nil {
		return userEnvVar, errors.Wrap(err, "could not match value")
	}

	if IsEnvVar {
		envVar, err := search(userEnvVar, `\w+`)
		if err != nil {
			return userEnvVar, errors.Wrap(err, "could not search value")
		}

		return os.Getenv(envVar), nil
	}

	return userEnvVar, nil
}

func match(str, pattern string) (bool, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, errors.Wrap(err, "failed to compile regex pattern")
	}

	return regex.MatchString(str), nil
}

func search(str, pattern string) (string, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", errors.Wrap(err, "failed to compile regex pattern")
	}

	return regex.FindString(str), nil
}
