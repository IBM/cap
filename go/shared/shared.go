/*
Copyright 2018 The cap Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shared

import (
	"time"
)

// NamedValue -
type NamedValue struct {
	ValueName string `xml:"valueName"`
	Value     string `xml:"value"`
}

// Search - returns the first value with name in Namevalue array or ""
func Search(nv *[]NamedValue, name string) string {
	for _, element := range *nv {
		if element.ValueName == name {
			return element.Value
		}
	}
	return ""
}

// SearchAll returns a slice of NamedValues for all values with the specified name
func SearchAll(nv *[]NamedValue, name string) []string {
	var found = make([]string, 0, len(*nv))

	for _, element := range *nv {
		if element.ValueName == name {
			found = append(found, element.Value)
		}
	}
	return found
}

// TimeStr - is a date/time in the CAPTimeFormat format
type TimeStr string

// Time - generate a new TimeStr from the passed in time.Time object
func Time(t time.Time) TimeStr {
	return TimeStr(t.Format(time.RFC3339))
}

// TimeParse - generate a time.Time from the passed in TimeStr
func TimeParse(t TimeStr) (time.Time, error) {
	return time.Parse(time.RFC3339, string(t))
}
