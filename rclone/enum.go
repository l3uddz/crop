package rclone

const (
	CmdCopy       string = "copy"
	CmdMove       string = "move"
	CmdSync       string = "sync"
	CmdDeleteFile string = "deletefile"
	CmdDeleteDir  string = "rmdir"
	CmdDeleteDirs string = "rmdirs"
	CmdDedupe     string = "dedupe"
)

const (
	ExitSuccess int = iota
	ExitSyntaxError
	ExitErrorUnknown
	ExitDirectoryNotFound
	ExitFileNotFound
	ExitTemporaryError
	ExitLessSeriousError
	ExitFatalError
	ExitTransferExceeded
)
