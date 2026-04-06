package domain

import (
	"fmt"
)

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

func (v Version) IncByCommit(c Commit) Version {
	if c.IsBreaking {
		v.IncMajor()
	} else if c.Type == "feat" {
		v.IncMinor()
	} else {
		v.IncPatch()
	}
	return v
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

type Commit struct {
	Raw        string
	Type       string
	Scope      string
	Message    string
	IsBreaking bool
}

func (c Commit) String() string {
	return c.Raw
}
