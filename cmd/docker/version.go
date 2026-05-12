package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/luvx21/coding-go/coding-common/sets"
)

type Version struct {
	Orinal       string
	Major        int
	Minor, Patch *int
}

var versionRegex = regexp.MustCompile(
	`(?P<major>0|[1-9][0-9]*)(?:\.(?P<minor>0|[1-9][0-9]*))?(?:\.(?P<patch>0|[1-9][0-9]*))?`,
)

// FromTag parses a tag string and returns Version and format string
// Example: "v4.0.3-ls215" -> Version{Major: 4, Minor: 0, Patch: 3}, "v{}.{}.{}-ls215"
func FromTag(tag string) (*Version, string, bool) {
	matches := versionRegex.FindAllStringSubmatchIndex(tag, -1)
	if len(matches) == 0 {
		return nil, "", false
	}

	// Find the best match with the most components
	var bestMatch []int
	maxComponents := 0

	for _, match := range matches {
		// Count matched groups (major, minor, patch)
		components := 0
		if match[2] >= 0 { // major matched
			components++
		}
		if match[4] >= 0 { // minor matched
			components++
		}
		if match[6] >= 0 { // patch matched
			components++
		}

		if components > maxComponents {
			maxComponents = components
			bestMatch = match
		}
	}

	if len(bestMatch) < 7 {
		return nil, "", false
	}

	// Extract major
	majorStr := tag[bestMatch[2]:bestMatch[3]]
	major, err := strconv.Atoi(majorStr)
	if err != nil {
		return nil, "", false
	}

	// Extract minor (optional)
	var minor *int
	if bestMatch[4] >= 0 {
		minorStr := tag[bestMatch[4]:bestMatch[5]]
		m, err := strconv.Atoi(minorStr)
		if err != nil {
			return nil, "", false
		}
		minor = &m
	}

	// Extract patch (optional)
	var patch *int
	if bestMatch[6] >= 0 {
		patchStr := tag[bestMatch[6]:bestMatch[7]]
		p, err := strconv.Atoi(patchStr)
		if err != nil {
			return nil, "", false
		}
		patch = &p
	}

	// Build format string
	positions := [][2]int{}
	if bestMatch[2] >= 0 {
		positions = append(positions, [2]int{bestMatch[2], bestMatch[3]})
	}
	if bestMatch[4] >= 0 {
		positions = append(positions, [2]int{bestMatch[4], bestMatch[5]})
	}
	if bestMatch[6] >= 0 {
		positions = append(positions, [2]int{bestMatch[6], bestMatch[7]})
	}

	formatStr := buildFormatString(tag, positions)

	return &Version{Orinal: tag, Major: major, Minor: minor, Patch: patch}, formatStr, true
}

// buildFormatString replaces matched positions with {}
func buildFormatString(tag string, positions [][2]int) string {
	// Sort positions in reverse order to replace from end to start
	for i := len(positions) - 1; i >= 0; i-- {
		pos := positions[i]
		tag = tag[:pos[0]] + "{}" + tag[pos[1]:]
	}
	return tag
}

func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	// Major versions are equal, compare minor
	switch {
	case v.Minor == nil && other.Minor == nil:
		return comparePatch(v.Patch, other.Patch)
	case v.Minor != nil && other.Minor != nil:
		if *v.Minor != *other.Minor {
			if *v.Minor < *other.Minor {
				return -1
			}
			return 1
		}
		return comparePatch(v.Patch, other.Patch)
	default:
		return 0
	}
}

// comparePatch compares two patch versions
func comparePatch(p1, p2 *int) int {
	switch {
	case p1 == nil && p2 == nil:
		return 0
	case p1 != nil && p2 != nil:
		if *p1 < *p2 {
			return -1
		} else if *p1 > *p2 {
			return 1
		}
		return 0
	default:
		return 0
	}
}
func formatCompatible(v1, v2 *Version, fmt1, fmt2 string) bool {
	if v1 == nil || v2 == nil {
		return false
	}
	switch {
	case (v1.Minor != nil && v2.Minor != nil) || (v1.Minor == nil && v2.Minor == nil):
		switch {
		case (v1.Patch != nil && v2.Patch != nil) || (v1.Patch == nil && v2.Patch == nil):
			return fmt1 == fmt2
		}
	}
	return false
}

func deduplicateVersions(versions []*Version) []*Version {
	seen := sets.NewSet[string]()
	var r []*Version
	for _, v := range versions {
		key := v.String()
		if !seen.Contains(key) {
			seen.Add(key)
			r = append(r, v)
		}
	}
	return r
}

func (v *Version) String() string {
	s := fmt.Sprintf("%d", v.Major)
	if v.Minor != nil {
		s += fmt.Sprintf(".%d", *v.Minor)
	}
	if v.Patch != nil {
		s += fmt.Sprintf(".%d", *v.Patch)
	}
	return s
}
