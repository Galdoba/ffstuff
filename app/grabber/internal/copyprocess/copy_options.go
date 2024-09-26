package copyprocess

type Option func(*copyOptions)

type copyOptions struct {
	sourcePaths   []string
	destination   string
	atemptLimit   int
	nameTabLen    int
	markerExt     string
	deleteMarkers bool
	deleteAll     bool
}

func defaultCopyOptions() copyOptions {
	return copyOptions{
		atemptLimit:   5,
		nameTabLen:    -1,
		deleteMarkers: true,
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

func WithMarkerExt(ext string) Option {
	return func(co *copyOptions) {
		co.markerExt = ext
	}
}
