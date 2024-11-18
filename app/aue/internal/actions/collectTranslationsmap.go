package actions

type translaitionProposal struct {
	inFileName     string
	inDirectory    string
	translatedName string
	words          []string
	score          float64
}

func CollectTranslationVariants(files ...string) {}
