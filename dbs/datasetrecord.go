package dbs

import lexicon "github.com/CHESSComputing/golib/lexicon"

// DatasetRecord represents input dataset record from HTTP request
type DatasetRecord struct {
	Did          string              `json:"did" validate:"required"`
	Buckets      []string            `json:"buckets" validate:"required"`
	Site         string              `json:"site" validate:"required"`
	Processing   string              `json:"processing" validate:"required"`
	Parent       string              `json:"parent_did" validate:"required"`
	InputFiles   []FileRecord        `json:"input_files,omitempty" validate:"required"`
	OutputFiles  []FileRecord        `json:"output_files,omitempty" validate:"required"`
	Environments []EnvironmentRecord `json:"environments"`
	Scripts      []ScriptRecord      `json:"scripts"`
	OsInfo       OsInfoRecord        `json:"osinfo"`
}

// Validate implementation of DatasetRecord
func (r *DatasetRecord) Validate() error {
	if err := lexicon.CheckPattern("did", r.Did); err != nil {
		return Error(err, ValidateErrorCode, "fail did validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("site", r.Site); err != nil {
		return Error(err, ValidateErrorCode, "fail site validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("processing", r.Processing); err != nil {
		return Error(err, ValidateErrorCode, "fail processing string validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("dataset_parent", r.Parent); err != nil {
		return Error(err, ValidateErrorCode, "fail dataset parent validation", "dbs.DatasetRecord.Validate")
	}
	for _, b := range r.Buckets {
		if err := lexicon.CheckPattern("bucket", b); err != nil {
			return Error(err, ValidateErrorCode, "fail bucket validation", "dbs.DatasetRecord.Validate")
		}
	}
	for _, f := range r.InputFiles {
		if err := lexicon.CheckPattern("fail", f.Name); err != nil {
			return Error(err, ValidateErrorCode, "fail file validation", "dbs.DatasetRecord.Validate")
		}
	}
	for _, f := range r.OutputFiles {
		if err := lexicon.CheckPattern("fail", f.Name); err != nil {
			return Error(err, ValidateErrorCode, "fail file validation", "dbs.DatasetRecord.Validate")
		}
	}
	if err := lexicon.CheckPattern("osinfo_name", r.OsInfo.Name); err != nil {
		return Error(err, ValidateErrorCode, "fail osinfo.name validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("osinfo_version", r.OsInfo.Version); err != nil {
		return Error(err, ValidateErrorCode, "fail osinfo.version validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("osinfo_kernel", r.OsInfo.Kernel); err != nil {
		return Error(err, ValidateErrorCode, "fail osinfo.name validation", "dbs.DatasetRecord.Validate")
	}
	for _, env := range r.Environments {
		if err := lexicon.CheckPattern("env_name", env.Name); err != nil {
			return Error(err, ValidateErrorCode, "fail env.Name validation", "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("env_version", env.Version); err != nil {
			return Error(err, ValidateErrorCode, "fail env.Version validation", "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("env_details", env.Details); err != nil {
			return Error(err, ValidateErrorCode, "fail env.Details validation", "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("env_parent", env.Parent); err != nil {
			return Error(err, ValidateErrorCode, "fail env.Parent validation", "dbs.DatasetRecord.Validate")
		}
	}
	for _, script := range r.Scripts {
		if err := lexicon.CheckPattern("script_name", script.Name); err != nil {
			return Error(err, ValidateErrorCode, "fail script.Name validation", "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("script_options", script.Options); err != nil {
			return Error(err, ValidateErrorCode, "fail script.Options validation", "dbs.DatasetRecord.Validate")
		}
		if err := lexicon.CheckPattern("script_parent", script.Parent); err != nil {
			return Error(err, ValidateErrorCode, "fail script.Parent validation", "dbs.DatasetRecord.Validate")
		}
	}
	return nil
}
