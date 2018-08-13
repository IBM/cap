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

package atom

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getNwsAtomFeedExample() (*Feed, error) {
	xmlData, err := ioutil.ReadFile("../../resources/nws_atom_feed_example.xml")
	if err != nil {
		return nil, err
	}
	var feed Feed
	err = xml.Unmarshal(xmlData, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func TestUnmarshalNWSAtomFeedHasProperValues(t *testing.T) {
	feed, err := getNwsAtomFeedExample()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, feed.ID, "https://alerts.weather.gov/cap/us.php?x=0")
	assert.Equal(t, feed.Logo, "http://alerts.weather.gov/images/xml_logo.gif")
	assert.Equal(t, feed.Generator.Content, "NWS CAP Server")
	assert.Equal(t, string(feed.Updated), "2018-08-15T16:57:00-06:00")
	assert.Equal(t, len(feed.Author), 1)
	assert.Equal(t, feed.Author[0].Name, "w-nws.webmaster@noaa.gov")
	assert.Equal(t, feed.Title.Content, "Current Watches, Warnings and Advisories for the United States Issued by the National Weather Service")
	assert.Equal(t, feed.Link[0].Href, "https://alerts.weather.gov/cap/us.php?x=0")
	assert.Equal(t, len(feed.Entries), 163)
}

func TestUnmarshalNWSAtomFeedEntryHasProperValues(t *testing.T) {
	feed, err := getNwsAtomFeedExample()
	if err != nil {
		t.Fatal(err)
	}
	var entry = feed.Entries[0]
	assert.Equal(t, entry.ID, "https://alerts.weather.gov/cap/wwacapget.php?x=AK125AB652A170.HighWindWarning.125AB660BDF0AK.AFGNPWNSB.e9d4afdcacb3b7015f58bccc1db60d46")
	assert.Equal(t, string(entry.Updated), "2018-08-15T14:52:00-08:00")
	assert.Equal(t, string(entry.Published), "2018-08-15T14:52:00-08:00")
	assert.Equal(t, len(entry.Author), 1)
	assert.Equal(t, entry.Author[0].Name, "w-nws.webmaster@noaa.gov")
	assert.Equal(t, entry.Title.Content, "High Wind Warning issued August 15 at 2:52PM AKDT until August 16 at 7:00AM AKDT by NWS")
	assert.Equal(t, entry.Link[0].Href, "https://alerts.weather.gov/cap/wwacapget.php?x=AK125AB652A170.HighWindWarning.125AB660BDF0AK.AFGNPWNSB.e9d4afdcacb3b7015f58bccc1db60d46")
	assert.Equal(t, entry.Summary.Content, "...HIGH WIND WARNING REMAINS IN EFFECT UNTIL 7 AM AKDT THURSDAY... * WINDS...Southwest 30 to 40 mph with gusts to 60 mph. * TIMING...Strong winds this evening will continue through Thursday morning. The strongest winds are expected late this evening. Winds will decrease early Thursday morning. * IMPACTS...Loose objects may be blown away.")
	assert.Equal(t, entry.Event, "High Wind Warning")
	assert.Equal(t, string(entry.Effective), "2018-08-15T14:52:00-08:00")
	assert.Equal(t, string(entry.Expires), "2018-08-16T07:00:00-08:00")
	assert.Equal(t, entry.Status, "Actual")
	assert.Equal(t, entry.MsgType, "Alert")
	assert.Equal(t, len(entry.Category), 1)
	assert.Equal(t, entry.Category[0].Content, "Met")
	assert.Equal(t, entry.Urgency, "Expected")
	assert.Equal(t, entry.Severity, "Severe")
	assert.Equal(t, entry.Certainty, "Likely")
	assert.Equal(t, entry.AreaDesc, "Eastern Beaufort Sea Coast")
	assert.Equal(t, entry.Polygon[0], "")
	assert.Equal(t, len(entry.Geocode.Names), 2)
	assert.Equal(t, len(entry.Geocode.Values), 2)
}

func TestUnmarshalNWSAtomFeedEntryGeocodeHasProperValues(t *testing.T) {
	feed, err := getNwsAtomFeedExample()
	if err != nil {
		t.Fatal(err)
	}
	var entry = feed.Entries[0]
	var geocode = entry.Geocode
	assert.Equal(t, "002185", geocode.GetGeocodes("FIPS6")[0])
	assert.Equal(t, "AKZ204", geocode.GetGeocodes("UGC")[0])
}

func TestUnmarshalNWSAtomFeedEntryParameterHasProperValues(t *testing.T) {
	feed, err := getNwsAtomFeedExample()
	if err != nil {
		t.Fatal(err)
	}
	var entry = feed.Entries[0]
	assert.Contains(t, entry.GetParameter("VTEC"), "/O.CON.PAFG.HW.W.0011.180816T0000Z-180816T1500Z/")
}

func TestNWSAtomGeocodeGetValuesReturnsEmptyArrIfNotFound(t *testing.T) {
	var geocode Geocode
	found := geocode.GetGeocodes("not-a-real-key")
	assert.Equal(t, len(found), 0)
}

func TestLinkFollowAlertReturnsErrorForInvalidURL(t *testing.T) {
	link := Link{Href: "abcdef"}
	_, err := link.GetAlert()
	assert.Equal(t, "Get abcdef: unsupported protocol scheme \"\"", err.Error())
}

func TestHandleHttpResponseReturnsStartingErr(t *testing.T) {
	existingError := errors.New("prexisting error")
	_, err := handleHTTPResponse(nil, existingError)
	assert.Equal(t, existingError, err)
}

func TestHandleHttpResponseReturnsErrOnNon200StatusCode(t *testing.T) {
	var response http.Response
	response.StatusCode = 400
	_, err := handleHTTPResponse(&response, nil)
	assert.Equal(t, "HTTP status code: 400", err.Error())
}

func TestHandleHttpResponseReturnsErrOnZeroContentLength(t *testing.T) {
	var response http.Response
	response.StatusCode = 200
	response.ContentLength = 0
	_, err := handleHTTPResponse(&response, nil)
	assert.Equal(t, "No content", err.Error())
}
