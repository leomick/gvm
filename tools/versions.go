package tools

import (
	"slices"
	"strconv"
	"strings"
)

func SortVersions(versions []string) ([]string, error) {
	var intVersions [][]int
	for _, v := range versions {
		sliceStringVer := strings.Split(v, ".")
		var intVersion []int
		for _, v := range sliceStringVer {
			val, err := strconv.Atoi(v)
			if err != nil {
				return []string{}, err
			}
			intVersion = append(intVersion, val)
		}
		intVersions = append(intVersions, intVersion)
	}
	slices.SortFunc(intVersions, func(a, b []int) int {
		var uselength int
		if len(a) < len(b) {
			uselength = len(a)
		} else {
			uselength = len(b)
		}
		for i := 0; i < uselength; i++ {
			if a[i] > b[i] {
				return -1
			}
			if a[i] < b[i] {
				return 1
			}
		}
		return 0
	})
	var sorted []string
	for _, v := range intVersions {
		var stringVersion []string
		for _, versionsection := range v {
			stringVersion = append(stringVersion, strconv.Itoa(versionsection))
		}
		versionString := strings.Join(stringVersion, ".")
		sorted = append(sorted, versionString)
	}
	return sorted, nil
}
