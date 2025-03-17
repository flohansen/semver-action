package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var verReg = regexp.MustCompile(`(v)?([0-9]+)\.([0-9]+)\.([0-9]+)`)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) IncMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
}

func (v *Version) IncMinor() {
	v.Minor++
	v.Patch = 0
}

func (v *Version) IncPatch() {
	v.Patch++
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func NewVersionFromString(str string) (*Version, error) {
	match := verReg.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil, errors.New("invalid version string")
	}

	major, err := strconv.Atoi(match[2])
	if err != nil {
		return nil, fmt.Errorf("could not parse major version: %w", err)
	}

	minor, err := strconv.Atoi(match[3])
	if err != nil {
		return nil, fmt.Errorf("could not parse minot version: %w", err)
	}

	patch, err := strconv.Atoi(match[4])
	if err != nil {
		return nil, fmt.Errorf("could not parse patch version: %w", err)
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

var convCommReg = regexp.MustCompile(`([a-zA-Z ]+)(!)?(\(.*\))?:(.*)`)

type Commit struct {
	Raw        string
	Type       string
	Scope      string
	Message    string
	IsBreaking bool
}

func (c *Commit) String() string {
	return c.Raw
}

func NewCommitFromString(str string) (*Commit, error) {
	match := convCommReg.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil, errors.New("invalid commit message")
	}

	c := &Commit{}
	c.Raw = str
	c.Type = strings.ToLower(match[1])
	c.IsBreaking = c.Type == "breaking change" || match[2] == "!"
	return c, nil
}
