package utils

import (
	"encoding/xml"
	"io"
)

type XmlMap map[string]string

type xmlData struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m XmlMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	start = xml.StartElement{
		Name: xml.Name{
			Local: "xml",
		},
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		_ = e.Encode(xmlData{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func (m XmlMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	for {
		var e xmlData

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(m)[e.XMLName.Local] = e.Value
	}
	return nil
}
