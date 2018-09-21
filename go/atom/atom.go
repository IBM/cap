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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/IBM/cap/go/cap"
	"github.com/IBM/cap/go/shared"
)

// NwsNationalAtomFeedURL is the URL for the NWS National Atom feed
const NwsNationalAtomFeedURL string = "https://alerts.weather.gov/cap/us.php?x=1"

// TODO consider adding enums
// TODO add json conversion

// Feed - root structure of an http://www.w3.org/2005/Atom:feed
type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	CommonAttributes
	// required elements
	ID      string   `xml:"id"`               // identifies the feed using a universally unique and permanent URI. Eg. domain name.
	Title   Text     `xml:"title"`            // contains a human readable title for the feed
	Updated TimeStr  `xml:"updated"`          // indicates the last time the feed was modified in a significant way
	Author  []Person `xml:"author,omitempty"` // names at least one author of the feed. A feed must contain at least one author element unless all of the entry elements contain at least one author element
	Link    []Link   `xml:"link"`             // identifies a related Web page
	// optional elements
	Category    []Category  `xml:"category,omitempty"`    // specifies the categories that the feed belongs to
	Contributor []Person    `xml:"contributor,omitempty"` // names the contributors to the feed
	Generator   Generator   `xml:"generator,omitempty"`   // identifies the software used to generate the feed, for debugging and other purposes
	Icon        string      `xml:"icon,omitempty"`        // identifies a small image which provides iconic visual identification for the feed
	Logo        string      `xml:"logo,omitempty"`        // identifies a larger image which provides visual identification for the feed
	Rights      Text        `xml:"rights,omitempty"`      // conveys information about rights, e.g. copyrights, held in and over the feed
	SubTitle    Text        `xml:"subtitle,omitempty"`    // contains a human-readable description or subtitle for the feed
	Entries     []Entry     `xml:"entry,omitempty"`       // the entries of the feed
	Extension   []Extension `xml:",any,omitempty"`        // Custom extensions
}

// CommonAttributes - this struct is for common atom attributes
type CommonAttributes struct {
	Base string `xml:"base,attr,omitempty"` // Base - if not empty must be a valid url.URL.
	Lang string `xml:"lang,attr,omitempty"`
}

// Text - this struct is for an xml element that may include chardata and
// a type attribute.
type Text struct {
	CommonAttributes
	// Content - the chardata between the open and close of the tag element
	Content string `xml:",chardata"`
	// Type - optional attribute which determines how the content is encoded (default="text")
	//	If type="text", then this element contains plain text
	// 	If type="html", then this element contains entity escaped html.
	//  If type="xhtml", then this element contains inline xhtml, wrapped in a div element.
	Type string `xml:"type,attr,omitempty"`
	Src  string `xml:"src,attr,omitempty"` // Src if present must be a valid url.URL.
	Body string `xml:",innerxml"`          // use Body for xhtml text
}

// Person - describe a person, corporation, or similar entity. It has one
// required element, name, and two optional elements: uri, email.
type Person struct {
	CommonAttributes
	Name      string      `xml:"name"`            // conveys a human-readable name for the person
	URI       string      `xml:"uri,omitempty"`   // contains a home page for the person
	Email     string      `xml:"email,omitempty"` // contains an email address for the person
	Extension []Extension `xml:",any,omitempty"`  // Custom extensions
}

// Link - is patterned after html's link element. It has one required attribute,
// href, and five optional attributes: rel, type, hreflang, title, and length.
// Patterned after https://www.w3.org/TR/1999/REC-html401-19991224/struct/links.html#h-12.3
type Link struct {
	CommonAttributes
	Href string `xml:"href,attr"` // is the URI of the referenced resource (typically a Web page)
	// rel contains a single link relationship type. It can be a full URI or one of the following predefined values (default=alternate):
	//    alternate: an alternate representation of the entry or feed, for example a permalink to the html version of the entry, or the front page of the weblog.
	//    enclosure: a related resource which is potentially large in size and might require special handling, for example an audio or video recording.
	//    related: an document related to the entry or feed.
	//    self: the feed itself.
	//    via: the source of the information provided in the entry.
	Rel      string `xml:"rel,attr,omitempty"`
	Type     string `xml:"type,attr,omitempty"`     // indicates the media type of the resource
	HrefLang string `xml:"hreflang,attr,omitempty"` // indicates the language of the referenced resource
	Title    Text   `xml:"title,attr,omitempty"`    // human readable information about the link, typically for display purposes
	Length   string `xml:"length,attr,omitempty"`   // the length of the resource, in bytes
}

// Category - has one required attribute, term, and two optional attributes, scheme and label
// Example: <category term="technology"/>
type Category struct {
	CommonAttributes
	Content string `xml:",chardata"`
	Term    string `xml:"term,attr"`             // identifies the category
	Scheme  string `xml:"scheme,attr,omitempty"` // identifies the categorization scheme via a URI
	Label   string `xml:"label,attr,omitempty"`  // provides a human-readable label for display
}

// Generator - Identifies the software used to generate the feed, for debugging
// and other purposes, uri and version attributes are optional.
type Generator struct {
	CommonAttributes
	Content string `xml:",chardata"`
	URI     string `xml:"uri,attr,omitempty"`
	Version string `xml:"version,attr,omitempty"`
}

// Entry - (root) structure of an http://www.w3.org/2005/Atom:feed>entry
type Entry struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom entry"`
	CommonAttributes
	// *** required elements ***
	// ID - entry/id = cap:alert:identifier identifies the entry using a universally
	// unique and permanent URI, two entries can have the same id if they
	// represent the same entry at different times
	ID string `xml:"id"`
	// Title - entry/title = cap:alert:info:headline human readable title for the entry
	Title Text `xml:"title"`
	// Updated - indicates the last time the entry was modified in a significant way
	Updated TimeStr `xml:"updated"`

	// *** recommended elements ***
	// Author names at least one author of the entry
	Author []Person `xml:"author"`
	// Content - entry/content: = (embedded cap <alert>)
	// contains or links to the complete content of the entry (must be provided if
	// there is no alternate link, and should be provided if there is no summary)
	Content Text `xml:"content,omitempty"`
	// Link - entry/link = (URL of full CAP alert)
	// Linking to where the full CAP alert is hosted is recommended. Note the
	// following example practices when linking to the CAP alert:
	// •	Use an absolute (not a relative) URL.
	// •	The alert page must exist when the alert link is added to the feed.
	// Otherwise, feed clients that try to load the alert will encounter errors
	// at load time.
	// •	The link to the CAP alert must be correctly identified
	// (type="application/cap+xml"). Similarly, if other links are provided to
	// content, they must be appropriately associated with the correct MIME type
	// and other attributes.
	Link []Link `xml:"link"`
	// Summary entry/summary = cap:alert:info:headline or excerpt from cap:alert:info:description
	// conveys a short summary, abstract, or excerpt of the entry (summary should
	// be provided if there either is no content provided for the entry, or that
	// content is not inline (i.e., contains a src attribute), or if the content
	// is encoded in base64)
	Summary Text `xml:"summary"`

	// *** optional elements ***
	// Category - category = cap:alert:info:category
	// specifies the categories that the entry belongs to
	Category []Category `xml:"category,omitempty"`
	// Contributor - names the contributors to the entry
	Contributor []Person `xml:"contributor,omitempty"`
	// Published - entry/published = cap:alert:sent
	// Note: The CAP sent element is specified as the date and time of the
	// origination of the CAP alert. You may wish to use the date and time of the
	// publication of the feed instead, since this is the expected value for
	// <published>.
	Published TimeStr  `xml:"published"`
	Rights    Text     `xml:"rights,omitempty"` // conveys information about rights, e.g. copyrights, held in and over the entry
	Source    []Source `xml:"source,omitempty"` // metadata from the source feed(s) if this entry is a copy
	// atom entry extensions (TODO all optional?) is there no spec for the entry extensions?
	// TODO research why they did not use Atom Extension for the additional fields and if this is kosher
	Event     string         `xml:"event,omitempty"`     // Event - The text denoting the type of the subject event in the alert entry.
	Effective shared.TimeStr `xml:"effective,omitempty"` // Effective - The effective date and time of the information in the alert entry.
	Expires   shared.TimeStr `xml:"expires,omitempty"`   // Expires - The expiry date and time of the information in the alert entry.
	Status    string         `xml:"status,omitempty"`    // Status - The code denoting the appropriate handling of the alert message.
	MsgType   string         `xml:"msgType,omitempty"`   // MsgType - The code denoting the nature of the alert message.
	Urgency   string         `xml:"urgency,omitempty"`   // Urgency - Urgency of the subject event of the alert entry.
	Severity  string         `xml:"severity,omitempty"`  // Severity - Severity of the subject event of the alert entry.
	Certainty string         `xml:"certainty,omitempty"` // Certainty - Certainty of the subject event of the alert entry.
	AreaDesc  string         `xml:"areaDesc,omitempty"`  // AreadDesc - The text describing the affected area of the alert entry.
	Polygon   []string       `xml:"polygon,omitempty"`   // Polygon - The paired values of points defining a polygon that delineates the affected area of the alert entry.
	Circle    []string       `xml:"circle,omitempty"`    // Circle - Note: in the xsd but not explained in CAP 1.2 documentation
	Geocode   Geocode        `xml:"geocode,omitempty"`   // Geocode - The geographic code delineating the affected area of the alert message.
	Parameter []NamedValue   `xml:"parameter,omitempty"` // Parameter - Denotes additional information associated with the alert entry.
	Extension []Extension    `xml:",any,omitempty"`      // Custom extensions
}

// Source - metadata from the source feed for entries that are a copy
type Source struct { // TODO if possible replace these struct elements with "Feed *feed"
	ID          string     `xml:"id"`                    // identifies the source of an entry using a universally unique and permanent URI
	Title       Text       `xml:"title"`                 // human readable title for the source of the entry
	Updated     string     `xml:"updated"`               // indicates the last time the source entry was modified in a significant way
	Author      []Person   `xml:"author,omitempty"`      // names at least one author of the source feed. A feed must contain at least one author element unless all of the entry elements contain at least one author element
	Link        []Link     `xml:"link"`                  // identifies a related Web page
	Category    []Category `xml:"category,omitempty"`    // specifies the categories that the source feed belongs to
	Contributor []Person   `xml:"contributor,omitempty"` // names the contributors to the source feed
	Generator   Generator  `xml:"generator,omitempty"`   // identifies the software used to generate the source feed, for debugging and other purposes
	Icon        string     `xml:"icon,omitempty"`        // identifies a small image which provides iconic visual identification for the source feed
	Logo        string     `xml:"logo,omitempty"`        // identifies a larger image which provides visual identification for the source feed
	Rights      Text       `xml:"rights,omitempty"`      // conveys information about rights, e.g. copyrights, held in and over the source feed
	SubTitle    Text       `xml:"subtitle,omitempty"`    // contains a human-readable description or subtitle for the source feed
}

// Geocode - as specified by the NWS Atom feed
type Geocode struct {
	Names  []string `xml:"valueName,omitempty"`
	Values []string `xml:"value,omitempty"`
}

// Extension - used for adding custom content to an element
type Extension struct {
	XMLName xml.Name
	XML     string `xml:",innerxml"`
}

// NamedValue -
type NamedValue = shared.NamedValue

// TimeStr - is a date/time in the AtomTimeFormat format
type TimeStr = shared.TimeStr

// TODO finish up the helper functions

// Time - generate a new TimeStr from the passed in time.Time object
func Time(t time.Time) TimeStr {
	return shared.Time(t)
}

// TimeParse - generate a time.Time from the passed in TimeStr
func TimeParse(t TimeStr) (time.Time, error) {
	return shared.TimeParse(t)
}

// GetGeocodes returns back an array of values for the Geocode element with the same name
func (g *Geocode) GetGeocodes(name string) []string {
	for index, value := range g.Names {
		if value == name {
			return strings.Split(g.Values[index], " ")
		}
	}
	return []string{}
}

// GetParameter returns the value for the first parameter with the specified name or ""
func (e *Entry) GetParameter(name string) string {
	return shared.Search(&e.Parameter, name)
}

// GetAlert retrieves an Alert from a link's href attribute
func (l *Link) GetAlert() (*cap.Alert11, []byte, error) {
	body, err := handleHTTPResponse(http.Get(l.Href))
	if err != nil {
		return nil, nil, err
	}

	var alert cap.Alert11
	err = xml.Unmarshal(body, &alert)
	if err != nil {
		return nil, nil, err
	}
	return &alert, body, nil
}

// GetFeed retrieves the main National Weather Service CAP v1.1 ATOM feed
func GetFeed() (*Feed, []byte, error) {
	return GetFeedFrom(NwsNationalAtomFeedURL)
}

// GetFeedFrom retrieves a CAP v1.1 ATOM feed from requested host
func GetFeedFrom(host string) (*Feed, []byte, error) {
	body, err := handleHTTPResponse(http.Get(host))
	if err != nil {
		return nil, nil, err
	}

	var downloadedFeed Feed
	err = xml.Unmarshal(body, &downloadedFeed)
	if err != nil {
		return nil, nil, err
	}
	return &downloadedFeed, body, nil
}

func handleHTTPResponse(r *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code: %d", r.StatusCode)
	}
	if r.ContentLength == 0 {
		return nil, fmt.Errorf("No content")
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
