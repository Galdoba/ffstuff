package config

const (
	IN          = "IN         "
	IN_PROGRESS = "IN_PROGRESS"
	DONE        = "DONE       "
	EDIT        = "EDIT       "
	ARCHIVE     = "ARCHIVE    "
)

type OptFunc func(*Options)

func defaultOptions() Options {
	pathsMap := make(map[string]string)
	pathsMap[IN] = ""
	pathsMap[IN_PROGRESS] = ""
	pathsMap[DONE] = ""
	pathsMap[EDIT] = ""
	pathsMap[ARCHIVE] = ""

	return Options{
		PATH:         pathsMap,
		CycleSeconds: 30,
	}
}

type Options struct {
	PATH            map[string]string `json:"Paths"`
	CycleSeconds    int               `json:"Repeat Cycle (Seconds)"`
	SkipRenameStage bool              `json:"Skip Rename Stage"`
	LogFilePath     string            `json:"Log file path,omitempty"` //later
}

func WithPath(key, path string) OptFunc {
	return func(opt *Options) {
		opt.PATH[key] = path
	}
}
