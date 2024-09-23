package actions

type Option func(*copyOptions)

type copyOptions struct {
	atemptLimit       int
	nameTabLen        int
	formatter         func(string, ...interface{}) string
	copyExistDecidion string
	copyPrefix        string
	copySuffix        string
}

func defaultCopyOptions() copyOptions {
	return copyOptions{
		atemptLimit: 5,
		nameTabLen:  -1,
	}
}
