package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type AmediaMetadataJson struct {
	Movies []struct {
		Actors         string `json:"actors"`
		AgeRestriction string `json:"age_restriction"`
		CmsID          int64  `json:"cms_id"`
		Country        string `json:"country"`
		Directors      string `json:"directors"`
		EndDate        string `json:"end_date"`
		File           struct {
			Duration float64 `json:"duration"`
			Serid    string  `json:"serid"`
		} `json:"file"`
		Genre          string `json:"genre"`
		GUID           string `json:"guid"`
		ImdbID         string `json:"imdb_id"`
		KinopoiskID    string `json:"kinopoisk_id"`
		OriginalTitle  string `json:"original_title"`
		Quote          string `json:"quote"`
		QuoteAuthor    string `json:"quote_author"`
		RusDescription string `json:"rus_description"`
		RusTitle       string `json:"rus_title"`
		StartDate      string `json:"start_date"`
		Year           int64  `json:"year"`
	} `json:"movies"`
	Series []struct {
		AgeRestriction      string `json:"age_restriction"`
		CmsID               int64  `json:"cms_id"`
		Country             string `json:"country"`
		EndDate             string `json:"end_date"`
		Genre               string `json:"genre"`
		GUID                string `json:"guid"`
		ImdbID              string `json:"imdb_id"`
		KinopoiskID         string `json:"kinopoisk_id"`
		OriginalBroadcaster string `json:"original_broadcaster"`
		OriginalTitle       string `json:"original_title"`
		Quote               string `json:"quote"`
		QuoteAuthor         string `json:"quote_author"`
		RusDescription      string `json:"rus_description"`
		RusTitle            string `json:"rus_title"`
		Seasons             []struct {
			Actors    string `json:"actors"`
			CmsID     int64  `json:"cms_id"`
			Directors string `json:"directors"`
			EndDate   string `json:"end_date"`
			Episodes  []struct {
				CmsID           int64  `json:"cms_id"`
				EndDate         string `json:"end_date"`
				EpisodeSynopsis string `json:"episode_synopsis"`
				File            struct {
					Duration float64 `json:"duration"`
					Serid    string  `json:"serid"`
				} `json:"file"`
				GUID                string `json:"guid"`
				OrderNumber         int64  `json:"order_number"`
				OriginalEpisodeName string `json:"original_episode_name"`
				RusEpisodeName      string `json:"rus_episode_name"`
				StartDate           string `json:"start_date"`
				Year                int64  `json:"year"`
			} `json:"episodes"`
			GUID              string `json:"guid"`
			OrderNumber       int64  `json:"order_number"`
			OrigName          string `json:"orig_name"`
			RusName           string `json:"rus_name"`
			SeasonDescription string `json:"season_description"`
			StartDate         string `json:"start_date"`
			YearsI            int64  `json:"years,omitempty"`
			YearsS            string `json:"years,omitempty"`
		} `json:"seasons"`
		StartDate string `json:"start_date"`
		YearsI    int64  `json:"years,omitempty"`
		YearsS    string `json:"years,omitempty"`
	} `json:"series"`
}

func TranslationsMap() (map[string]string, error) {
	characters = map[string]int{}
	trMap := make(map[string]string)
	meta := AmediaMetadataJson{}
	bt, err := os.ReadFile(`\\192.168.31.4\buffer\IN\@AMEDIA_IN\metadata.json`)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	if err := json.Unmarshal(bt, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %v", err)
	}

	for _, movie := range meta.Movies {
		asciiKey := normalizeText(movie.OriginalTitle)
		asciiValue := normalizeText(movie.RusTitle)
		renameTo := translation(words(asciiValue))
		if renameTo == "?" {
			renameTo = translation(words(asciiKey))
		}
		trMap[asciiKey] = renameTo
	}
	for _, serial := range meta.Series {
		asciiKey := normalizeText(serial.OriginalTitle)
		asciiValue := normalizeText(serial.RusTitle)
		renameTo := translation(words(asciiValue))
		if renameTo == "?" {
			renameTo = translation(words(asciiKey))
		}
		trMap[asciiKey] = renameTo
	}

	return trMap, err
}

func inDirectoryWords(dir string) []string {
	//Best_Interests_105_PRT241118000535
	dir = strings.TrimSuffix(dir, prtStrSuffix(dir))
	dir = strings.TrimSuffix(dir, seNumSuffix(dir))
	return words(strings.ToLower(dir))
}

func prtStrSuffix(str string) string {
	re := regexp.MustCompile(`(_PRT[0-9]{12,})$`)
	return re.FindString(str)
}

func seNumSuffix(str string) string {
	re := regexp.MustCompile(`(_[0-9]{3,})$`)
	return re.FindString(str)
}

func guessFolder(name string) string {
	for _, glyph := range []string{":", "-", "_", ",", "&", ".", ";", "#", "'", `"`, "?", "/", `\`, "|", "(", ")", "’", "*", "!", "[", "]", "{", "}"} {
		name = strings.ReplaceAll(name, glyph, " ")
	}
	words := strings.Fields(name)
	name = strings.Join(words, "_")
	return name
}

var characters map[string]int

func count(s string) {
	for _, ch := range strings.Split(s, "") {
		characters[ch]++
	}
}

func hasNonLatin(s string) bool {
	s = strings.ToLower(s)
	for _, ch := range strings.Split(s, "") {
		switch ch {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "_", " ":
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
		case "а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и", "й", "к", "л", "м", "н", "о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь", "э", "ю", "я":
		case "'", ".", ":", "?", "&", ",", "`", " ", "-", "*", "(", ")", "/", "➔", "!", "#", "’":
		case "½", "è", "é":
		default:

			return true
		}

	}

	return false
}

func normalizeText(text string) string {
	text = strings.ToLower(text)
	normalized := ""
	for _, ch := range strings.Split(text, "") {
		switch ch {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
		case "а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и", "й", "к", "л", "м", "н", "о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь", "э", "ю", "я":
		case "-", ".", "_", " ", "'", ":", "?", "&", ",", "`", " ", "*", "(", ")", "/", "➔", "!", "#", "’", "½", "°", "—", `"`, "«":
			ch = "_"
		case "è", "é":
			ch = "e"
		default:
			ch = "-"
		}
		normalized += ch

	}
	return normalized
}

func words(text string) []string {
	text = strings.ReplaceAll(text, "_", " ")
	return strings.Fields(text)
}

func translation(words []string) string {
	text := strings.Join(words, "_")
	transliterated := ""
	lMap := map[string]string{
		"а": "a", "б": "b", "в": "v", "г": "g", "д": "d",
		"е": "e", "ё": "e", "ж": "zh", "з": "z", "и": "i",
		"й": "y", "к": "k", "л": "l", "м": "m", "н": "n",
		"о": "o", "п": "p", "р": "r", "с": "s", "т": "t",
		"у": "u", "ф": "f", "х": "h", "ц": "c", "ч": "ch",
		"ш": "sh", "щ": "sh", "ъ": "", "ы": "y", "ь": "",
		"э": "e", "ю": "yu", "я": "ya", "_": "_"}
	for _, ch := range strings.Split(text, "") {
		switch ch {
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
			"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
			"u", "v", "w", "x", "y", "z",
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			transliterated += ch
		default:
			transliterated += lMap[ch]
		}

	}
	if len(transliterated) < 1 {
		return "?"
	}
	sl := strings.Split(transliterated, "")[0]
	sl = strings.ToUpper(sl)
	return sl + strings.Join(strings.Split(transliterated, "")[1:], "")
}
