package template

import "github.com/tomdoesdev/knox/pkg/errs"

const (
	TemplateParseFailureCode   errs.ErrorCode = "TEMPLATE_PARSE_FAILURE"
	TemplateExecuteFailureCode errs.ErrorCode = "TEMPLATE_EXECUTE_FAILURE"
	TemplateFileNotFoundCode   errs.ErrorCode = "TEMPLATE_FILE_NOT_FOUND"
	TemplateSecretNotFoundCode errs.ErrorCode = "TEMPLATE_SECRET_NOT_FOUND"
	TemplateReadFailureCode    errs.ErrorCode = "TEMPLATE_READ_FAILURE"
)

var (
	ErrTemplateFileNotFound = errs.New(TemplateFileNotFoundCode, "template file not found")
)
