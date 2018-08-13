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
	"encoding/xml"
	"time"

	"github.com/IBM/cap/go/shared"
)

// TODO consider adding enums
// TODO add json conversion

// Alert - This struct is for a CAP Alert Message (version 1.2)
type Alert struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:emergency:cap:1.2 alert"`

	Identifier  string   `xml:"identifier"`            // Identifier - A string which uniquely identifies the CAP message.
	Sender      string   `xml:"sender"`                // Sender - Email address of the NWS webmaster.
	Sent        TimeStr  `xml:"sent"`                  // Sent - The origination time and date of the alert message.
	Status      string   `xml:"status"`                // Status - The code denoting the appropriate handling of the alert message.
	MsgType     string   `xml:"msgType"`               // MsgType - The code denoting the nature of the alert message.
	Source      string   `xml:"source,omitempty"`      // Source - Note: in the xsd but not explained in CAP 1.2 documentation
	Scope       string   `xml:"scope"`                 // Scope - The code denoting the appropriate handling of the alert message.
	Restriction string   `xml:"restriction,omitempty"` // Restriction: Note: in the xsd but not explained in CAP 1.2 documentation
	Addresses   string   `xml:"addresses,omitempty"`   // Addresses - Note: in the xsd but not explained in CAP 1.2 documentation
	Code        []string `xml:"code,omitempty"`        // Code - Version of the CAP IPAWS profile as adopted by FEMA to which the subject CAP message conforms.
	Note        string   `xml:"note,omitempty"`        // Note - The text describing the purpose or significance of the alert message.
	References  []string `xml:"references,omitempty"`  // References - References the most recent message to which the current message refers or replaces.
	Incidents   []string `xml:"incidents,omitempty"`   // Incidents - Note: in the xsd but not explained in CAP 1.2 documentation
	Info        []Info   `xml:"info,omitempty"`        // Info - The container for all component parts of the info element.
}

// Alert11 CAP v1.1 Alert Message
type Alert11 struct {
	Alert
	XMLName xml.Name `xml:"urn:oasis:names:tc:emergency:cap:1.1 alert"` // TODO ensure this is actually a duplicate
}

// Info -
type Info struct {
	XMLName xml.Name `xml:"info"`

	Language     string       `xml:"language,omitempty"`     // Language - Note: language is specified in the CAP xsd but not in the CAP v1.2 documentation, for details on use see http://www.datypic.com/sc/xsd/t-xsd_language.html
	Category     []string     `xml:"category"`               // Category - The code denoting the category of the subject event in the alert message. Multiple instances may occur within an <info> block.
	Event        string       `xml:"event"`                  // Event - The text denoting the type of the subject event in the alert message
	ResponseType []string     `xml:"responseType,omitempty"` // ResponseType - The code denoting the type of action recommended for the target audience.
	Urgency      string       `xml:"urgency"`                // Urgency - Urgency of the subject event of the alert message.
	Severity     string       `xml:"severity"`               // Severity - Severity of the subject event of the alert message.
	Certainty    string       `xml:"certainty"`              // Certainty - Certainty of the subject event of the alert message.
	Audience     string       `xml:"audience,omitempty"`     // Audience - is in the CAP xsd but not in the CAP v1.2 documentation
	EventCode    []NamedValue `xml:"eventCode,omitempty"`    // EventCode - A system-specific code identifying the event type of the alert message.
	Effective    TimeStr      `xml:"effective,omitempty"`    // Effective - The effective date and time of the information in the alert message.
	Onset        TimeStr      `xml:"onset,omitempty"`        // Onset - Expected time of the beginning of the subject event in the alert message.
	Expires      TimeStr      `xml:"expires,omitempty"`      // Expires - The expiry date and time of the information in the alert message.
	SenderName   string       `xml:"senderName,omitempty"`   // SenderName - Name of the issuing NWS Office.
	Headline     string       `xml:"headline,omitempty"`     // Headline - A brief human-readable headline containing the alert type and valid time of the alert.
	Description  string       `xml:"description,omitempty"`  // Description - The text describing the subject event of the alert message.
	Instruction  string       `xml:"instruction,omitempty"`  // Instruction - The text describing the recommended action to be taken by recipients of the alert message.
	Web          string       `xml:"web,omitempty"`          // Web - A hyperlink where additional information about the alert can be found.
	Contact      string       `xml:"contact,omitempty"`      // Contact - Note: in the xsd but not explained in CAP 1.2 documentation
	Parameter    []NamedValue `xml:"parameter,omitempty"`    // Parameter - Denotes additional information associated with the alert message.
	Resource     []Resource   `xml:"resource,omitempty"`     // Resource - in the xsd but not explained in CAP 1.2 documentation.
	Area         []Area       `xml:"area,omitempty"`         // Area - array of area elements associated with the alert message.
}

// Resource - Note: in the xsd but not explained in CAP 1.2 documentation
type Resource struct {
	XMLName xml.Name `xml:"resource"`

	ResourceDesc string `xml:"resourceDesc"`
	MIMEType     string `xml:"mimeType"`
	Size         int64  `xml:"size,omitempty"`
	URI          string `xml:"uri,omitempty"`
	DerefURI     string `xml:"derefUri,omitempty"`
	Digest       string `xml:"digest,omitempty"`
}

// Area - The container for all sub-elements of the area element.
type Area struct {
	XMLName xml.Name `xml:"area"`

	AreaDesc string       `xml:"areaDesc"`           // AreadDesc - The text describing the affected area of the alert message.
	Polygon  []string     `xml:"polygon,omitempty"`  // Polygon - The paired values of points defining a polygon that delineates the affected area of the alert message.
	Circle   []string     `xml:"circle,omitempty"`   // Circle - Note: in the xsd but not explained in CAP 1.2 documentation
	Geocode  []NamedValue `xml:"geocode,omitempty"`  // Geocode - The geographic code delineating the affected area of the alert message.
	Altitude string       `xml:"altitude,omitempty"` // TODO need a xs:decimal type here // Note: in the xsd but not explained in CAP 1.2 documentation
	Ceiling  string       `xml:"ceiling,omitempty"`  // TODO need a xs:decimal type here // Note: in the xsd but not explained in CAP 1.2 documentation
}

// NamedValue -
type NamedValue = shared.NamedValue

// TimeStr - is a date/time in the CAPTimeFormat format
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

// ParseAlert parses XML bytes into a CAP 1.2 Alert
func ParseAlert(xmlData []byte) (*Alert, error) {
	var alert Alert

	err := xml.Unmarshal(xmlData, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// ParseAlert11 parses XML bytes into a CAP 1.1 Alert
func ParseAlert11(xmlData []byte) (*Alert11, error) {
	var alert Alert11

	err := xml.Unmarshal(xmlData, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// GetParameter returns back the value for the first parameter with the specified name
func (info *Info) GetParameter(name string) string {
	return shared.Search(&info.Parameter, name)
}

// AddParameter adds a Parameter with the specified name and value
func (info *Info) AddParameter(name string, value string) {
	param := NamedValue{ValueName: name, Value: value}
	info.Parameter = append(info.Parameter, param)
}

// GetGeocode returns back the value for the first Geocode value with the specified name
func (a *Area) GetGeocode(name string) string {
	return shared.Search(&a.Geocode, name)
}

// GetGeocodes returns back the Geocode values with the specified name
func (a *Area) GetGeocodes(name string) []string {
	return shared.SearchAll(&a.Geocode, name)
}

// AddGeocode adds a Geocode with the specified name and value
func (a *Area) AddGeocode(name string, value string) {
	geocode := NamedValue{ValueName: name, Value: value}
	a.Geocode = append(a.Geocode, geocode)
}
