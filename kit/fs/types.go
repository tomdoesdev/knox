package fs

import "github.com/tomdoesdev/knox/kit/errs"

const (
	ECodeFileReadFailure  errs.Code = "FILE_READ_FAILURE"
	ECodeFileWriteFailure errs.Code = "FILE_WRITE_FAILURE"
	ECodeFileMoveFailure  errs.Code = "FILE_MOVE_FAILURE"

	ECodeDirectoryFailure errs.Code = "DIRECTORY_FAILURE"

	ECodeTempFailure errs.Code = "TEMP_FAILURE"

	ECodeInvalidPath errs.Code = "INVALID_PATH"

	ECodeEntityExists errs.Code = "ENTITY_EXISTS"
)
