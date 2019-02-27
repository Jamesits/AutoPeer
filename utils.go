package main

import (
	"reflect"
	"regexp"
	"strings"
)

// https://codereview.stackexchange.com/a/60085
func inArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

var uidCounter = 0

func getUid() int {
	uidCounter += 1
	return uidCounter
}

func cleanString(s string) string {
	s = strings.Trim(s, " ")
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)

	rep := re.ReplaceAllString(s, "_")

	// prevent first character being a digit
	return "AP_" + rep
}
