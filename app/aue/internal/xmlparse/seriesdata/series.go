package seriesdata

import (
	"encoding/xml"
	"fmt"
)

type VideoData struct {
	Group []struct {
		Guid  string `xml:"guid,attr"`
		Type  string `xml:"type,attr"`
		Group []struct {
			Number   int    `xml:"number,attr"`
			Type     string `xml:"type,attr"`
			MetaInfo struct {
				Available struct {
					Start    string `xml:"start,attr"`
					CharData string `xml:",chardata"`
				} `xml:"available"`
				Description *struct {
					Type     string `xml:"type,attr"`
					CharData string `xml:",chardata"`
				} `xml:"description"`
				Title struct {
					Type     string `xml:"type,attr"`
					CharData string `xml:",chardata"`
				} `xml:"title"`
				Year int `xml:"year"`
			} `xml:"meta-info"`
			Video []struct {
				End             int    `xml:"end,attr"`
				Endtitles       int    `xml:"endtitles,attr"`
				Episodesinopsys string `xml:"episodesinopsys,attr"`
				Guid            string `xml:"guid,attr"`
				Multilang       bool   `xml:"multilang,attr"`
				Number          int    `xml:"number,attr"`
				Src             string `xml:"src,attr"`
				Start           int    `xml:"start,attr"`
				Logo            struct {
					Src      string `xml:"src,attr"`
					CharData string `xml:",chardata"`
				} `xml:"logo"`
				MetaInfo struct {
					Available struct {
						End   *string `xml:"end,attr"`
						Start string  `xml:"start,attr"`
					} `xml:"available"`
					Duration int  `xml:"duration"`
					Featured *int `xml:"featured"`
					Title    []struct {
						Type     string `xml:"type,attr"`
						CharData string `xml:",chardata"`
					} `xml:"title"`
				} `xml:"meta-info"`
				Subtitles struct {
					Src      string `xml:"src,attr"`
					CharData string `xml:",chardata"`
				} `xml:"subtitles"`
			} `xml:"video"`
		} `xml:"group"`
		MetaInfo struct {
			Available bool   `xml:"available"`
			Category  string `xml:"category"`
			Credits   struct {
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
			Featured    string  `xml:"featured"`
			ImdbID      *string `xml:"imdb_id"`
			KinopoiskID *string `xml:"kinopoisk_id"`
			Location    string  `xml:"location"`
			Priority    int     `xml:"priority"`
			Quote       *struct {
				Author   string `xml:"author,attr"`
				CharData string `xml:",chardata"`
			} `xml:"quote"`
			Restriction string `xml:"restriction"`
			Slogan      *struct {
				Type     string `xml:"type,attr"`
				CharData string `xml:",chardata"`
			} `xml:"slogan"`
			StudioRestrictions *struct {
				EpisodesAllowed int `xml:"episodes_allowed"`
			} `xml:"studio_restrictions"`
			Title []struct {
				Type     string `xml:"type,attr"`
				CharData string `xml:",chardata"`
			} `xml:"title"`
			Year int `xml:"year"`
		} `xml:"meta-info"`
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
