package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	commitExp  = regexp.MustCompile(`(?<type>[a-zA-Z ]+)(\((?<scope>.*)\))?(?<breaking>!)?:(?<message>.*)`)
	versionExp = regexp.MustCompile(`(v)?(?<major>[0-9]+)\.(?<minor>[0-9]+)\.(?<patch>[0-9]+)`)
)

func DecodeCommitFromString(str string) (Commit, error) {
	matches := commitExp.FindStringSubmatch(str)
	if len(matches) == 0 {
		return Commit{}, errors.New("invalid commit message")
	}

	commit := Commit{
		Raw: str,
	}

	if i := commitExp.SubexpIndex("type"); i >= 0 && i < len(matches) {
		commit.Type = matches[i]
		if commit.Type == "breaking change" {
			commit.IsBreaking = true
		}
	}
	if i := commitExp.SubexpIndex("breaking"); i >= 0 && i < len(matches) {
		commit.IsBreaking = matches[i] == "!"
	}
	if i := commitExp.SubexpIndex("scope"); i >= 0 && i < len(matches) {
		commit.Scope = matches[i]
	}
	if i := commitExp.SubexpIndex("message"); i >= 0 && i < len(matches) {
		commit.Message = strings.TrimSpace(matches[i])
	}

	return commit, nil
}

func DecodeVersionFromString(str string) (Version, error) {
	matches := versionExp.FindStringSubmatch(str)
	if len(matches) == 0 {
		return Version{}, errors.New("invalid version string")
	}

	version := Version{}

	if i := versionExp.SubexpIndex("major"); i >= 0 && i < len(matches) {
		major, err := strconv.Atoi(matches[i])
		if err != nil {
			return Version{}, fmt.Errorf("could not parse major version: %w", err)
		}
		version.Major = major
	}

	if i := versionExp.SubexpIndex("minor"); i >= 0 && i < len(matches) {
		minor, err := strconv.Atoi(matches[i])
		if err != nil {
			return Version{}, fmt.Errorf("could not parse minor version: %w", err)
		}
		version.Minor = minor
	}

	if i := versionExp.SubexpIndex("patch"); i >= 0 && i < len(matches) {
		patch, err := strconv.Atoi(matches[i])
		if err != nil {
			return Version{}, fmt.Errorf("could not parse patch version: %w", err)
		}
		version.Patch = patch
	}

	return version, nil
}
