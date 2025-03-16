package dbs

import (
	"fmt"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// DatasetRecord represents input dataset record from HTTP request
type DatasetRecord struct {
	Did          string              `json:"did" validate:"required"`
	Site         string              `json:"site" validate:"required"`
	Processing   string              `json:"processing" validate:"required"`
	Parent       string              `json:"parent_did" validate:"required"`
	InputFiles   []FileRecord        `json:"input_files,omitempty" validate:"required"`
	OutputFiles  []FileRecord        `json:"output_files,omitempty" validate:"required"`
	Environments []EnvironmentRecord `json:"environments"`
	Scripts      []ScriptRecord      `json:"scripts"`
	Buckets      []BucketRecord      `json:"buckets"`
	OsInfo       OsInfoRecord        `json:"osinfo"`
}

// Validate implementation of DatasetRecord
func (r *DatasetRecord) Validate() error {
	if err := lexicon.CheckPattern("did", r.Did); err != nil {
		msg := fmt.Sprintf("fail did validation: '%v'", r.Did)
		return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("site", r.Site); err != nil {
		msg := fmt.Sprintf("fail site validation: '%v'", r.Site)
		return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("processing", r.Processing); err != nil {
		msg := fmt.Sprintf("fail processing validation: '%v'", r.Processing)
		return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("dataset_parent", r.Parent); err != nil {
		msg := fmt.Sprintf("fail parent validation: '%v'", r.Parent)
		return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
	}
	for _, b := range r.Buckets {
		if err := lexicon.CheckPattern("bucket", b.Name); err != nil {
			msg := fmt.Sprintf("fail bucket validation: '%v'", b.Name)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
	}
	for _, f := range r.InputFiles {
		if err := lexicon.CheckPattern("file", f.Name); err != nil {
			msg := fmt.Sprintf("fail file validation: '%v'", f.Name)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
	}
	for _, f := range r.OutputFiles {
		if err := lexicon.CheckPattern("fail", f.Name); err != nil {
			msg := fmt.Sprintf("fail file validation: '%v'", f.Name)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
	}
	if err := lexicon.CheckPattern("osinfo_name", r.OsInfo.Name); err != nil {
		msg := fmt.Sprintf("fail osinfo validation: '%v'", r.OsInfo.Name)
		return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("osinfo_version", r.OsInfo.Version); err != nil {
		msg := fmt.Sprintf("fail osinfo validation: '%v'", r.OsInfo.Version)
		return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("osinfo_kernel", r.OsInfo.Kernel); err != nil {
		msg := fmt.Sprintf("fail osinfo validation: '%v'", r.OsInfo.Kernel)
		return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
	}
	for _, env := range r.Environments {
		if err := lexicon.CheckPattern("env_name", env.Name); err != nil {
			msg := fmt.Sprintf("fail env_name validation: '%v'", env.Name)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("env_version", env.Version); err != nil {
			msg := fmt.Sprintf("fail env_version validation: '%v'", env.Version)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("env_details", env.Details); err != nil {
			msg := fmt.Sprintf("fail env_details validation: '%v'", env.Details)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("env_parent", env.Parent); err != nil {
			msg := fmt.Sprintf("fail env_parent validation: '%v'", env.Parent)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
	}
	for _, script := range r.Scripts {
		if err := lexicon.CheckPattern("script_name", script.Name); err != nil {
			msg := fmt.Sprintf("fail script_name validation: '%v'", script.Name)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("script_options", script.Options); err != nil {
			msg := fmt.Sprintf("fail script_options validation: '%v'", script.Options)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("script_parent", script.Parent); err != nil {
			msg := fmt.Sprintf("fail script_parent validation: '%v'", script.Parent)
			return Error(err, ValidateErrorCode, msg, "dbs.DatasetRecord.Validate")
		}
	}
	return nil
}
