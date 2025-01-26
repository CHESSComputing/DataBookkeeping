package dbs

// DBS errors module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"errors"
	"fmt"
	"runtime"
)

// GenericErr represents generic dbs error
var GenericErr = errors.New("dbs error")

// DatabaseErr represents generic database error
var DatabaseErr = errors.New("database error")

// InvalidParamErr represents generic error for invalid input parameter
var InvalidParamErr = errors.New("invalid parameter(s)")

// ConcurrencyErr represents generic concurrency error
var ConcurrencyErr = errors.New("concurrency error")

// RecordErr represents generic record error
var RecordErr = errors.New("record error")

// ValidationErr represents generic validation error
var ValidationErr = errors.New("validation error")

// ContentTypeErr represents generic content-type error
var ContentTypeErr = errors.New("content-type error")

// NotImplementedApiErr represents generic not implemented api error
var NotImplementedApiErr = errors.New("not implemented api error")

// InvalidRequestErr represents generic invalid request error
var InvalidRequestErr = errors.New("invalid request error")

// DBS Error codes provides static representation of DBS errors, they cover 1xx range
const (
	GenericErrorCode        = iota + 100 // generic DBS error
	DatabaseErrorCode                    // 101 database error
	TransactionErrorCode                 // 102 transaction error
	QueryErrorCode                       // 103 query error
	RowsScanErrorCode                    // 104 row scan error
	SessionErrorCode                     // 105 db session error
	CommitErrorCode                      // 106 db commit error
	ParseErrorCode                       // 107 parser error
	LoadErrorCode                        // 108 loading error, e.g. load template
	GetIDErrorCode                       // 109 get id db error
	InsertErrorCode                      // 110 db insert error
	UpdateErrorCode                      // 111 update error
	DeleteErrorCode                      // 112 error error
	LastInsertErrorCode                  // 113 db last insert error
	ValidateErrorCode                    // 114 validation error
	PatternErrorCode                     // 115 pattern error
	DecodeErrorCode                      // 116 decode error
	EncodeErrorCode                      // 117 encode error
	ContentTypeErrorCode                 // 118 content type error
	ParametersErrorCode                  // 119 parameters error
	NotImplementedApiCode                // 120 not implemented API error
	ReaderErrorCode                      // 121 io reader error
	WriterErrorCode                      // 122 io writer error
	UnmarshalErrorCode                   // 123 json unmarshal error
	MarshalErrorCode                     // 124 marshal error
	HttpRequestErrorCode                 // 125 HTTP request error
	MigrationErrorCode                   // 126 Migration error
	RemoveErrorCode                      // 127 remove error
	InvalidRequestErrorCode              // 128 invalid request error
	PackagesErrorCode                    // 129 packages error code
	ManyToManyErrorCode                  // 130 many-to-many insertion error code
	DatasetErrorCode                     // 131 dataset error code
	EnvironmentsErrorCode                // 132 environments error code
	FilesErrorCode                       // 133 files error code
	OsInfoErrorCode                      // 134 osinfo error code
	ParentsErrorCode                     // 135 parents error code
	ProcessingErrorCode                  // 136 processing error code
	ScriptsErrorCode                     // 137 scripts error code
	SitesErrorCode                       // 138 sites error code
	LastAvailableErrorCode               // last available DBS error code
)

// DBSError represents common structure for DBS errors
type DBSError struct {
	Reason     string `json:"reason"`     // error string
	Message    string `json:"message"`    // additional message describing the issue
	Function   string `json:"function"`   // DBS function
	Code       int    `json:"code"`       // DBS error code
	Stacktrace string `json:"stacktrace"` // Go stack trace
}

// Error function implements details of DBS error message
func (e *DBSError) Error() string {
	return fmt.Sprintf(
		"\nDBSError\n   Code: %d\n   Description: %s\n   Function: %s\n   Message: %s\n   Reason: %v\n",
		e.Code, e.Explain(), e.Function, e.Message, e.Reason)
}

// ErrorStacktrace function implements details of DBS error message and stacktrace
func (e *DBSError) ErrorStacktrace() string {
	return fmt.Sprintf(
		"\nDBSError Stacktrace\n   Code: %d\n   Description: %s\n   Function: %s\n   Message: %s\n   Reason: %v\nStacktrace: %v\n\n",
		e.Code, e.Explain(), e.Function, e.Message, e.Reason, e.Stacktrace)
}

func (e *DBSError) Explain() string {
	switch e.Code {
	case GenericErrorCode:
		return "Generic DBS error"
	case DatabaseErrorCode:
		return "DBS DB error"
	case TransactionErrorCode:
		return "DBS DB transaction error"
	case QueryErrorCode:
		return "DBS DB query error, e.g. malformed SQL statement"
	case RowsScanErrorCode:
		return "DBS DB row scane error, e.g. fail to get DB record from a database"
	case SessionErrorCode:
		return "DBS DB session error"
	case CommitErrorCode:
		return "DBS DB transaction commit error"
	case ParseErrorCode:
		return "DBS parser error, e.g. malformed input parameter to the query"
	case LoadErrorCode:
		return "DBS file load error, e.g. fail to load DB template"
	case GetIDErrorCode:
		return "DBS DB ID error for provided entity, e.g. there is no record in DB for provided value"
	case InsertErrorCode:
		return "DBS DB insert record error"
	case UpdateErrorCode:
		return "DBS DB update record error"
	case LastInsertErrorCode:
		return "DBS DB laster insert record error, e.g. fail to obtain last inserted ID"
	case ValidateErrorCode:
		return "DBS validation error, e.g. input parameter does not match lexicon rules"
	case PatternErrorCode:
		return "DBS validation error when wrong pattern is provided"
	case DecodeErrorCode:
		return "DBS decode record failure, e.g. malformed JSON"
	case EncodeErrorCode:
		return "DBS encode record failure, e.g. unable to convert structure to JSON"
	case ContentTypeErrorCode:
		return "Wrong Content-Type HTTP header in HTTP request"
	case ParametersErrorCode:
		return "DBS invalid parameter for the DBS API"
	case NotImplementedApiCode:
		return "DBS Not implemented API error"
	case ReaderErrorCode:
		return "DBS reader I/O error, e.g. unable to read HTTP POST payload"
	case WriterErrorCode:
		return "DBS writer I/O error, e.g. unable to write record to HTTP response"
	case UnmarshalErrorCode:
		return "DBS unable to parse JSON record"
	case MarshalErrorCode:
		return "DBS unable to convert record to JSON"
	case HttpRequestErrorCode:
		return "invalid HTTP request"
	case RemoveErrorCode:
		return "Unable to remove record from DB"
	case InvalidRequestErrorCode:
		return "Invalid HTTP request"
	default:
		return "Not defined"
	}
	return "Not defined"
}

// helper function to create dbs error
func Error(err error, code int, msg, function string) error {
	reason := "nil"
	if err != nil {
		reason = err.Error()
	}
	stackSlice := make([]byte, 1024*4)
	s := runtime.Stack(stackSlice, false)
	return &DBSError{
		Reason:     reason,
		Message:    msg,
		Code:       code,
		Function:   function,
		Stacktrace: fmt.Sprintf("\n%s", stackSlice[0:s]),
	}
}
