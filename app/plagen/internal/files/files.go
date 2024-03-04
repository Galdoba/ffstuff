package files

import "os"

func ListDir(dir string) []string {
	fls, _ := os.ReadDir(dir)
	names := []string{}
	for _, fl := range fls {
		if fl.IsDir() {
			continue
		}
		names = append(names, dir+fl.Name())
	}
	return names
}
