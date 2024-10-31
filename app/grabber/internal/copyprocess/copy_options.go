package copyprocess

import (
	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
)

type Option func(*copyOptions)

type copyOptions struct {
	sourcePaths   []origin.Origin
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

func WithSourcePaths(sourcePaths ...origin.Origin) Option {
	return func(co *copyOptions) {
		for _, arg := range sourcePaths {
			if sliceContains(co.sourcePaths, arg) {
				continue
			}
			co.sourcePaths = append(co.sourcePaths, arg)
		}
	}
}

func sliceContains(sl []origin.Origin, s origin.Origin) bool {
	for _, have := range sl {
		if have.Path() == s.Path() {
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

func WithDeleteDecidion(d string) Option {
	return func(co *copyOptions) {
		switch d {
		case grabberflag.VALUE_DELETE_ALL:
			co.deleteAll = true
			co.deleteMarkers = true
		case grabberflag.VALUE_DELETE_MARKER:
			co.deleteMarkers = true
		}
	}
}
