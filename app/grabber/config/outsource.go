package config

func (cfg *Configuration) NamesPriority() map[string]int {
	return cfg.FILE_PRIORITY_WEIGHTS
}

func (cfg *Configuration) DirectoryPriority() map[string]int {
	return cfg.DIRECTORY_PRIORITY_WEIGHTS
}
