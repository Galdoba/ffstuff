package targetfile

type TargetFile struct {
	claimedGoal         string
	usedSourcesNames    []string
	expectedName        string
	expectedStreamTypes map[string]int
}

func New(goal string, sources []string) TargetFile {
	tf := TargetFile{}
	tf.claimedGoal = goal
	tf.usedSourcesNames = sources

	return tf
}
