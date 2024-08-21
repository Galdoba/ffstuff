package job

type JobOptsFunc func(*jobOptions)

type jobOptions struct {
	processingMode string
	jobType        string
}

func defaultJobOptions() jobOptions {
	return jobOptions{
		processingMode: "Undefined",
		jobType:        "Undefined",
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
