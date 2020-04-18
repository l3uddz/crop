package rclone

const (
	CMD_COPY        string = "copy"
	CMD_MOVE        string = "move"
	CMD_SYNC        string = "sync"
	CMD_DELETE_FILE string = "deletefile"
	CMD_DELETE_DIR  string = "rmdir"
	CMD_DELETE_DIRS string = "rmdirs"
	CMD_DEDUPE      string = "dedupe"
)

const (
	EXIT_SUCCESS int = iota
	EXIT_SYNTAX_ERROR
	EXIT_ERROR_UNKNOWN
	EXIT_DIRECTORY_NOT_FOUND
	EXIT_FILE_NOT_FOUND
	EXIT_TEMPORARY_ERROR
	EXIT_LESS_SERIOUS_ERROR
	EXIT_FATAL_ERROR
	EXIT_TRANSFER_EXCEEDED
)
