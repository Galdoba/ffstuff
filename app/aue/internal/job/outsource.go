package job

func (ja *jobAdmin) ProjectName() string {
	return ja.name
}

func (ja *jobAdmin) BashDestination() string {
	return ja.options.bashDestination
}

func (ja *jobAdmin) BashTranslationMap() map[string]string {
	return ja.options.bashTranslationMap
}
