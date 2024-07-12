package scan

type Instuction struct {
	Key string
	Val string
}

const (
	Key_TYPE   = "Check_Type"
	Key_FILE   = "File_Path"
	RDWR_Check = "Read-Write Check"
)

var RWCHECK Instuction = Instuction{
	Key: Key_TYPE,
	Val: RDWR_Check,
}

func File(path string) Instuction {
	return Instuction{
		Key: Key_FILE,
		Val: path,
	}
}
