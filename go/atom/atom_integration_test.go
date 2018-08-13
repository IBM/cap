// +build integration

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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNwsAtomFeed(t *testing.T) {
	feed, raw, err := GetFeed()
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 {
		t.Fatalf("Feed did not download correctly: %v", feed.ID)
	}
	if feed.ID != "https://alerts.weather.gov/cap/us.php?x=0" {
		t.Fatalf("Feed did not download / parse correctly: %v", feed.ID)
	}
}

func TestGetNwsAtomFeedEntryFromLink(t *testing.T) {
	feed, raw, err := GetFeed()
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 {
		t.Fatalf("Feed did not download correctly: %v", feed.ID)
	}
	if len(feed.Entries) == 0 {
		t.Skip("No alert entries in the Atom feed. Skipping...")
	}
	entry := feed.Entries[0]
	alert, err := entry.Link[0].GetAlert()
	if err != nil {
		t.Fatal(err)
	}
	entryID, err := url.Parse(entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "NOAA-NWS-ALERTS-"+entryID.Query().Get("x"), alert.Identifier)
}
