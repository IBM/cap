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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeParseReturnsCorrectValue(t *testing.T) {
	var nt TimeStr = "2003-06-11T22:39:00+01:00"
	dt, _ := TimeParse(nt)
	_, zoneOffset := dt.Zone()
	zoneOffsetHours := zoneOffset / 3600
	assert.Equal(t, 2003, int(dt.Year()))
	assert.Equal(t, 6, int(dt.Month()))
	assert.Equal(t, 11, int(dt.Day()))
	assert.Equal(t, 22, int(dt.Hour()))
	assert.Equal(t, 39, int(dt.Minute()))
	assert.Equal(t, 0, int(dt.Second()))
	assert.Equal(t, +1, zoneOffsetHours)
}
