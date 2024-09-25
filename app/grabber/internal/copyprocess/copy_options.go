package copyprocess

type Option func(*copyOptions)

type copyOptions struct {
	sourcePaths            []string
	destination            string
	atemptLimit            int
	nameTabLen             int
	formatter              func(string, ...interface{}) string
	copyExistDecidion      string
	copyPrefix             string
	copySuffix             string
	delete_original_source bool
	delete_original_marker bool
}

func defaultCopyOptions() copyOptions {
	return copyOptions{
		atemptLimit: 5,
		nameTabLen:  -1,
	}
}

func WithSourcePaths(sourcePaths ...string) Option {
	return func(co *copyOptions) {
		for _, arg := range sourcePaths {
			if sliceContains(co.sourcePaths, arg) {
				continue
			}
			co.sourcePaths = append(co.sourcePaths, arg)
		}
	}
}

func sliceContains(sl []string, s string) bool {
	for _, have := range sl {
		if have == s {
			return true
		}
	}
	return false
}

func WithDestination(dest string) Option {
	return func(co *copyOptions) {
		co.destination = dest
	}
}
