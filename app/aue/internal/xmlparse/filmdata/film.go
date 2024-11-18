package filmdata

import (
	"encoding/xml"
	"fmt"
)

type VideoData struct {
	Group struct {
		Type  string `xml:"type,attr"`
		Video []struct {
			End      int    `xml:"end,attr"`
			Guid     string `xml:"guid,attr"`
			Src      string `xml:"src,attr"`
			Start    int    `xml:"start,attr"`
			MetaInfo struct {
				Available struct {
					Start    string `xml:"start,attr"`
					CharData string `xml:",chardata"`
				} `xml:"available"`
				Category string `xml:"category"`
				Credits  struct {
					Credit []struct {
						Role     string `xml:"role,attr"`
						CharData string `xml:",chardata"`
						Award    []struct {
							Type     string `xml:"type,attr"`
							Year     int    `xml:"year,attr"`
							CharData string `xml:",chardata"`
						} `xml:"award"`
					} `xml:"credit"`
				} `xml:"credits"`
				Description struct {
					Type     string `xml:"type,attr"`
					CharData string `xml:",chardata"`
				} `xml:"description"`
				Duration    int     `xml:"duration"`
				Featured    string  `xml:"featured"`
				ImdbID      *string `xml:"imdb_id"`
				KinopoiskID *string `xml:"kinopoisk_id"`
				Location    string  `xml:"location"`
				Priority    int     `xml:"priority"`
				Restriction string  `xml:"restriction"`
				Slogan      *struct {
					Type     string `xml:"type,attr"`
					CharData string `xml:",chardata"`
				} `xml:"slogan"`
				Title []struct {
					Type     string `xml:"type,attr"`
					CharData string `xml:",chardata"`
				} `xml:"title"`
				Year int `xml:"year"`
			} `xml:"meta-info"`
		} `xml:"video"`
	} `xml:"group"`
	Title string `xml:"title"`
}

func Unmarshal(bt []byte) (*VideoData, error) {
	vid := VideoData{}
	err := xml.Unmarshal(bt, &vid)
	if err != nil {

		return nil, fmt.Errorf("filmfile unmarshaling failed: %v", err)
	}
	return &vid, nil
}
