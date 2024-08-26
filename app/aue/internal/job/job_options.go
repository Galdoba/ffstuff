package job

type JobOptsFunc func(*jobOptions)

type jobOptions struct {
	processingMode string
	jobType        string
	inputDir       string
	processingDir  string
	doneDir        string
	outDir         string
	bashGeneration bool
}

func defaultJobOptions() jobOptions {
	return jobOptions{
		processingMode: "",
		jobType:        "",
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

func WithBashGeneration(bg bool) JobOptsFunc {
	return func(opts *jobOptions) {
		opts.bashGeneration = bg
	}
}
