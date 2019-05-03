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

package commands

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"time"

	"github.com/IBM/cap/go/atom"
	"github.com/IBM/cap/go/cap"
	"github.com/IBM/cap/go/shared"
)

// generates an atom feed from a directory of alerts
// TODO feedGenerator aggregation from other feeds/feed entries
// TODO one feed per language? We've got alernative languages in the cap alerts to handle
func (s *Server) feedGenerater(path string) (*atom.Feed, error) {
	var entries []atom.Entry

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		xmlData, err := ioutil.ReadFile(path + file.Name())
		if err != nil {
			return nil, err
		}
		alert, err := cap.ParseAlert(xmlData)
		if err != nil {
			return nil, err
		}
		var categories []atom.Category
		for _, c := range alert.Info[0].Category {
			categories = append(categories, atom.Category{Content: c})
		}
		var (
			areas    string
			polygons []string
			circles  []string
			geocode  atom.Geocode
		)
		for _, a := range alert.Info[0].Area {
			areas = areas + "; " + a.AreaDesc // TODO semi colon is a guess
			polygons = append(polygons, a.Polygon...)
			circles = append(circles, a.Circle...)
			for _, g := range a.Geocode {
				geocode.Names = append(geocode.Names, g.ValueName)
				geocode.Values = append(geocode.Values, g.Value)
			}
		}
		entry := atom.Entry{
			ID: alert.Identifier,
			Title: atom.Text{
				Content: alert.Info[0].Headline, // TODO add dates? what to do here
			},
			Updated: alert.Sent, // TODO is this the right date?
			Author: []atom.Person{
				{
					Name: alert.Sender,
				},
			},
			//Content: ,
			Link: []atom.Link{
				{
					Href: s.config.HostName + alertsDirectory + file.Name(),
				},
			},
			Summary: atom.Text{
				Content: alert.Info[0].Description, // TODO hmm is this right?
			},
			Category: categories,
			// Contributor: []atom.Person `xml:"contributor,omitempty"`
			Published: alert.Sent,
			// Rights TODO how to get these
			// Source: *** Note atom.Source is only used for entry copies from other atom feeds
			Event:     alert.Info[0].Event,
			Effective: alert.Info[0].Effective,
			Expires:   alert.Info[0].Expires,
			Status:    alert.Status,
			MsgType:   alert.MsgType,
			Urgency:   alert.Info[0].Urgency,
			Severity:  alert.Info[0].Severity,
			Certainty: alert.Info[0].Certainty,
			AreaDesc:  areas,
			Polygon:   polygons,
			Circle:    circles,
			Geocode:   geocode,
		}
		entries = append(entries, entry)
	}

	// TODO get or configure these fields and the rest of the optional ones not initialized
	feed := &atom.Feed{
		ID: s.config.HostName + "cap/",
		Title: atom.Text{
			Content: "Current Alerts Issued by foo.com",
		},
		Updated: shared.TimeStr(time.Now().Format(time.RFC3339)),
		Author: []atom.Person{
			{
				Name: "foo.webmaster@foo.com", // TODO
			},
		},
		Link: []atom.Link{
			{
				Href: s.config.HostName + "cap/",
			},
		},
		Generator: atom.Generator{
			Content: "foo CAP Server",
		},
		Logo:    "http://alerts.weather.gov/images/xml_logo.gif",
		Entries: entries,
	}

	return feed, nil
}

// stores an atom feed to the feeds directory
// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (s *Server) feedWriter(feed *atom.Feed) (int, error) {
	b, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return 0, err
	}
	feedPath := s.config.Root + feedsDirectory + usFeed
	if err = os.Remove(feedPath); err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	feedFile, err := os.OpenFile(feedPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer feedFile.Close()
	return feedFile.Write(b)
}
