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

package cap

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getCAPAlertExample() (*Alert, error) {
	xmlData, err := ioutil.ReadFile("../../resources/cap_amber_alert_example.xml")
	if err != nil {
		return nil, err
	}
	return ParseAlert(xmlData)
}

func TestUnmarshalAlertHasProperValues(t *testing.T) {
	alert, err := getCAPAlertExample()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, alert.Identifier, "KAR0-0306112239-SW")
	assert.Equal(t, alert.Sender, "KARO@CLETS.DOJ.CA.GOV")
	assert.Equal(t, string(alert.Sent), "2003-06-11T22:39:00-07:00")
	assert.Equal(t, alert.Status, "Actual")
	assert.Equal(t, alert.MsgType, "Alert")
	assert.Equal(t, alert.Scope, "Public")
	assert.Equal(t, alert.Note, "")
	assert.Equal(t, len(alert.Info), 2)
}

func TestUnmarshalAlertInfoHasProperValues(t *testing.T) {
	alert, err := getCAPAlertExample()
	if err != nil {
		t.Fatal(err)
	}
	var info = alert.Info[0]
	assert.Equal(t, info.Category[0], "Rescue")
	assert.Equal(t, info.Event, "Child Abduction")
	assert.Equal(t, info.Urgency, "Immediate")
	assert.Equal(t, info.Certainty, "Likely")
	assert.Equal(t, info.EventCode[0].ValueName, "SAME")
	assert.Equal(t, info.EventCode[0].Value, "CAE")
	assert.Equal(t, string(info.Effective), "")
	assert.Equal(t, string(info.Expires), "")
	assert.Equal(t, info.SenderName, "Los Angeles Police Dept - LAPD")
	assert.Equal(t, info.Headline, "Amber Alert in Los Angeles County")
	assert.Contains(t, info.Description, "DATE/TIME: 06/11/03, 1915 HRS.  VICTIM(S): KHAYRI D")
	assert.Contains(t, info.Instruction, "")
	assert.Equal(t, len(info.Parameter), 0)
	assert.Equal(t, len(info.Area), 1)
}

func TestUnmarshalAlertInfoParameterHasProperValues(t *testing.T) {
	alert, err := getCAPAlertExample()
	if err != nil {
		t.Fatal(err)
	}
	var info = alert.Info[0]
	assert.Equal(t, info.GetParameter("WMOHEADER"), "")
	assert.Equal(t, info.GetParameter("TIME"), "")
}

func TestUnmarshalAlertInfoAreaHasProperValues(t *testing.T) {
	alert, err := getCAPAlertExample()
	if err != nil {
		t.Fatal(err)
	}
	var info = alert.Info[0]
	var area = info.Area[0]
	assert.Equal(t, area.AreaDesc, "Los Angeles County")
	assert.Equal(t, len(area.Polygon), 0)
	assert.Equal(t, len(area.Geocode), 1)
	assert.Equal(t, "006037", area.GetGeocodes("SAME")[0])
}

func TestAddParameterToInfoSetsProperValue(t *testing.T) {
	parameterName := "testcode"
	parameterValue := "1234"
	var info Info
	assert.Equal(t, len(info.Parameter), 0)
	info.AddParameter(parameterName, parameterValue)
	assert.Equal(t, len(info.Parameter), 1)
	parameter := info.Parameter[0]
	assert.Equal(t, parameter.ValueName, parameterName)
	assert.Equal(t, parameter.Value, parameterValue)
}

func TestAddGeocodeToAreaSetsProperValue(t *testing.T) {
	geocodeName := "testcode"
	geocodeValue := "1234"
	var area Area
	assert.Equal(t, len(area.Geocode), 0)
	area.AddGeocode(geocodeName, geocodeValue)
	assert.Equal(t, len(area.Geocode), 1)
	geocode := area.Geocode[0]
	assert.Equal(t, geocode.ValueName, geocodeName)
	assert.Equal(t, geocode.Value, geocodeValue)
}

func TestAreaGecodeReturnsFirstValue(t *testing.T) {
	geocode1 := NamedValue{ValueName: "test-name", Value: "1234"}
	geocode2 := NamedValue{ValueName: "test-name", Value: "5678"}
	var area Area
	area.AddGeocode(geocode1.ValueName, geocode1.Value)
	area.AddGeocode(geocode2.ValueName, geocode2.Value)
	assert.Equal(t, len(area.Geocode), 2)
	geocodeValue := area.GetGeocode("test-name")
	assert.Equal(t, geocodeValue, geocode1.Value)
}

func TestAreaGecodeReturnsEmptyStringIfNotFound(t *testing.T) {
	geocode := NamedValue{ValueName: "test-name", Value: "1234"}
	var area Area
	area.AddGeocode(geocode.ValueName, geocode.Value)
	geocodeValue := area.GetGeocode("not-a-real-key")
	assert.Equal(t, geocodeValue, "")
}

func TestParseCAPDateReturnsCorrectValue(t *testing.T) {
	alert, err := getCAPAlertExample()
	if err != nil {
		t.Fatal(err)
	}
	dt, _ := TimeParse(alert.Sent)
	_, zoneOffset := dt.Zone()
	zoneOffsetHours := zoneOffset / 3600
	// 2003-06-11T22:39:00-07:00
	assert.Equal(t, 2003, int(dt.Year()))
	assert.Equal(t, 6, int(dt.Month()))
	assert.Equal(t, 11, int(dt.Day()))
	assert.Equal(t, 22, int(dt.Hour()))
	assert.Equal(t, 39, int(dt.Minute()))
	assert.Equal(t, 0, int(dt.Second()))
	assert.Equal(t, -7, zoneOffsetHours)
}

func TestParseAlertReturnsErrForInvalidXml(t *testing.T) {
	_, err := ParseAlert([]byte("invalid xml"))
	assert.Equal(t, "EOF", err.Error())
}

func TestParseAlert11ReturnsErrForInvalidXml(t *testing.T) {
	_, err := ParseAlert11([]byte("invalid xml"))
	assert.Equal(t, "EOF", err.Error())
}
