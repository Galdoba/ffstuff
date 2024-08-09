package commands

import (
	"fmt"
	"regexp"
	"strings"
)

func sortAmediaFilesByEpisodes() (map[string][]string, error) {
	groups := make(map[string][]string)

	files, err := listAmediaFiles(in_dir)
	if err != nil {
		return groups, err
	}
	for _, file := range files {
		re := regexp.MustCompile(`(--s[0-9]{1,}e[0-9]{1,}--)`)
		episode := re.FindString(file)
		episode = strings.ReplaceAll(episode, "--", "")
		episode = strings.Split(file, "--")[0] + "--SER--" + episode
		groups[episode] = append(groups[episode], file)
	}
	if _, ok := groups[""]; ok {
		return groups, fmt.Errorf("un asighed group: %v", strings.Join(groups[""], " && "))
	}
	return groups, nil
}
