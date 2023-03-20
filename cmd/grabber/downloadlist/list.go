package downloadlist

type downloadStatus struct {
	source      string
	destination string
	transfered  int64
	complete    bool
}
