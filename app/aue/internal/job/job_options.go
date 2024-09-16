package job

type JobOptsFunc func(*jobOptions)

type jobOptions struct {
	processingMode     string
	jobType            string
	inputDir           string
	processingDir      string
	doneDir            string
	outDir             string
	outDirPrefix       string
	notificationDir    string
	directProcessing   bool
	bashGeneration     bool
	bashDestination    string
	bashTranslationMap map[string]string
}

func defaultJobOptions() jobOptions {
	return jobOptions{
		processingMode:     "",
		jobType:            "",
		bashTranslationMap: make(map[string]string),
	}
}

func WithProcessingMode(mode string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.processingMode = mode
	}
}

func WithJobType(jType string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.jobType = jType
	}
}

func WithInputDir(dir string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.inputDir = dir
	}
}

func WithProcessingDir(dir string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.processingDir = dir
	}
}

func WithDoneDir(dir string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.doneDir = dir
	}
}

func WithOutDir(dir string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.outDir = dir
	}
}

func WithNotificationDir(dir string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.notificationDir = dir
	}
}

func WithDirectProcessing(dp bool) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.directProcessing = dp
	}
}

func WithBashGeneration(bg bool) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.bashGeneration = bg
	}
}

func WithBashDestination(dest string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.bashDestination = dest
	}
}

func WithBashTranslationMap(trMap map[string]string) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.bashTranslationMap = trMap
	}
}
