package base

import (
	"fmt"
	"text/template"
)

// generalized error types

// basic errors
type errBasicError struct {
	Description, Suggestion string
}

func basicErrorString(d, s string) string {
	if s != "" {
		return fmt.Sprintf("license: %s\nlicense: %s", d, s)
	} else {
		return fmt.Sprintf("license: %s", d)
	}
}

// data errors
type errDataError struct {
	Description, Suggestion string
	Data                    interface{}
}

func dataErrorString(d, s string, data interface{}) string {
	if s != "" {
		return fmt.Sprintf("license: %s %v\nlicense: %s", d, data, s)
	} else {
		return fmt.Sprintf("license: %s %v", d, data)
	}
}

// argument errors
type errArgumentError struct {
	Description, Suggestion string
	Args                    []string
}

func argumentErrorString(d, s string, args []string) string {
	if s != "" {
		return fmt.Sprintf("license: %s %v\nlicense: %s", d, args, s)
	} else {
		return fmt.Sprintf("license: %s %v", d, args)
	}
}

// path errors
type errPathError struct {
	Description, Suggestion string
	Paths                   []string
}

func pathErrorString(d, s string, paths []string) string {
	if s != "" {
		return fmt.Sprintf("license: %s %v\nlicense: %s", d, paths, s)
	} else {
		return fmt.Sprintf("license: %s %v", d, paths)
	}
}

// specific error types

// basic errors
type errReadFailed errBasicError
type errFetchFailed errBasicError
type errParsingArguments errBasicError
type errCannotLocateHomeDir errBasicError
type errExpectedLicenseName errBasicError
type errCannotFindLicense errBasicError

func (err *errReadFailed) Error() string {
	return basicErrorString(err.Description, err.Suggestion)
}
func (err *errFetchFailed) Error() string {
	return basicErrorString(err.Description, err.Suggestion)
}
func (err *errParsingArguments) Error() string {
	return basicErrorString(err.Description, err.Suggestion)
}
func (err *errCannotLocateHomeDir) Error() string {
	return basicErrorString(err.Description, err.Suggestion)
}
func (err *errExpectedLicenseName) Error() string {
	return basicErrorString(err.Description, err.Suggestion)
}
func (err *errCannotFindLicense) Error() string {
	return basicErrorString(err.Description, err.Suggestion)
}

// data errors
type errSerializeFailed errDataError
type errDeserializeFailed errDataError
type errLoadingTemplate errDataError
type errExecutingTemplate errDataError

func (err *errSerializeFailed) Error() string {
	return dataErrorString(err.Description, err.Suggestion, err.Data)
}
func (err *errDeserializeFailed) Error() string {
	return dataErrorString(err.Description, err.Suggestion, err.Data)
}
func (err *errLoadingTemplate) Error() string {
	return dataErrorString(err.Description, err.Suggestion, err.Data)
}
func (err *errExecutingTemplate) Error() string {
	return dataErrorString(err.Description, err.Suggestion, err.Data)
}

// argument errors
type errUnknownArgument errArgumentError
type errBadArgumentSyntax errArgumentError

func (err *errUnknownArgument) Error() string {
	return argumentErrorString(err.Description, err.Suggestion, err.Args)
}
func (err *errBadArgumentSyntax) Error() string {
	return argumentErrorString(err.Description, err.Suggestion, err.Args)
}

// path errors
type errCreateTempDirFailed errPathError
type errWriteFileFailed errPathError
type errCreateDirFailed errPathError
type errRemovePathFailed errPathError

func (err *errCreateTempDirFailed) Error() string {
	return pathErrorString(err.Description, err.Suggestion, err.Paths)
}
func (err *errWriteFileFailed) Error() string {
	return pathErrorString(err.Description, err.Suggestion, err.Paths)
}
func (err *errCreateDirFailed) Error() string {
	return pathErrorString(err.Description, err.Suggestion, err.Paths)
}
func (err *errRemovePathFailed) Error() string {
	return pathErrorString(err.Description, err.Suggestion, err.Paths)
}

// copy tree error
type errCopyTreeFailed struct {
	From, To string
}

func (err *errCopyTreeFailed) Error() string {
	return fmt.Sprintf("license: failed to copy tree from %s to %s", err.From, err.To)
}

// constructors

// basic errors
func NewErrReadFailed() error {
	return &errReadFailed{
		"failed to read license(s)",
		"try again after running \"license update\"",
	}
}

func NewErrFetchFailed() error {
	return &errFetchFailed{
		"failed to fetch license(s)",
		"check your internet connection and try again",
	}
}

func NewErrParsingArguments() error {
	return &errParsingArguments{
		"error parsing arguments", "",
	}
}

func NewErrCannotLocateHomeDir() error {
	return &errCannotLocateHomeDir{
		"unable to locate home directory", "",
	}
}

func NewErrExpectedLicenseName() error {
	return &errExpectedLicenseName{
		"expected license name",
		"see \"license help\" for more details",
	}
}

func NewErrCannotFindLicense() error {
	return &errCannotFindLicense{
		"unable to find given command or license",
		"run \"license ls\" for a list of available licenses or see \"license help\"",
	}
}

// data errors
func NewErrSerializeFailed(l interface{}) error {
	return &errSerializeFailed{
		"failed to serialize license(s)",
		fmt.Sprintf("please create an issue at %s", repositoryIssuesUrl),
		l,
	}
}

func NewErrDeserializeFailed(data []byte) error {
	return &errDeserializeFailed{
		"failed to deserialize license(s)",
		fmt.Sprintf("please create an issue at %s", repositoryIssuesUrl),
		string(data),
	}
}

func NewErrLoadingTemplate(name string) error {
	return &errLoadingTemplate{
		"failed to load template",
		fmt.Sprintf("please create an issue at %s", repositoryIssuesUrl),
		name,
	}
}

func NewErrExecutingTemplate(t *template.Template) error {
	return &errExecutingTemplate{
		"failed to execute template",
		fmt.Sprintf("please create an issue at %s", repositoryIssuesUrl),
		t,
	}
}

// path errors
func NewErrCreateTempDirFailed(p ...string) error {
	return &errCreateTempDirFailed{
		"failed to create temporary directory", "", p,
	}
}

func NewErrCreateDirFailed(p ...string) error {
	return &errCreateDirFailed{
		"failed to create directory", "", p,
	}
}

func NewErrRemovePathFailed(p ...string) error {
	return &errRemovePathFailed{
		"failed to remove path", "", p,
	}
}

func NewErrWriteFileFailed(p ...string) error {
	return &errWriteFileFailed{
		"failed to write file", "", p,
	}
}

// argument errors
func NewErrUnknownArgument(args ...string) error {
	return &errUnknownArgument{
		"unknown argument",
		"see \"license help\" for more details",
		args,
	}
}

func NewErrBadFlagSyntax(args ...string) error {
	return &errBadArgumentSyntax{
		"bad flag",
		"see \"license help\" for more details",
		args,
	}
}

// copy tree error
func NewErrCopyTreeFailed(from, to string) error {
	return &errCopyTreeFailed{From: from, To: to}
}
