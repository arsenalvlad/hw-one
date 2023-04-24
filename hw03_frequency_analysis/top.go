package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type resStruct struct {
	key   string
	value int64
}

func Top10(str string) []string {
	res := make([]string, 0, 10)
	strMap := make(map[string]int64)
	reg1 := regexp.MustCompile(`([a-zA-Z]+[,.!'?:])*`)
	symbol := ",.!'?:"

	strSlice := strings.Fields(strings.ToLower(str))

	for _, k := range strSlice {
		if k == "-" {
			continue
		}
		if reg1.MatchString(k) {
			for _, l := range symbol {
				k = strings.ReplaceAll(k, string(l), "")
			}
		}
		if strMap[k] == 0 {
			strMap[k] = 1
		} else {
			strMap[k]++
		}
	}

	resSort := make([]resStruct, 0, len(strMap))
	for i, m := range strMap {
		resSort = append(resSort, resStruct{i, m})
	}

	sort.Slice(resSort, func(i, j int) bool {
		if resSort[i].value == resSort[j].value {
			return resSort[i].key < resSort[j].key
		}

		return resSort[i].value > resSort[j].value
	})

	for i, k := range resSort {
		if i > 9 {
			return res
		}
		res = append(res, k.key)
	}

	return res
}
