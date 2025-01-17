package dbs

import lexicon "github.com/CHESSComputing/golib/lexicon"

// DatasetRecord represents input dataset record from HTTP request
type DatasetRecord struct {
	Did         string            `json:"did" validate:"required"`
	Buckets     []string          `json:"buckets" validate:"required"`
	Site        string            `json:"site" validate:"required"`
	Processing  string            `json:"processing" validate:"required"`
	Parent      string            `json:"parent" validate:"required"`
	InputFiles  []string          `json:"input_files" validate:"required"`
	OutputFiles []string          `json:"output_files" validate:"required"`
	Environment EnvironmentRecord `json:"environment",omitempty`
	OsInfo      OsInfoRecord      `json:"osinfo",omitempty`
	Script      ScriptRecord      `json:"script",omitempty`
}

// Validate implementation of DatasetRecord
func (r *DatasetRecord) Validate() error {
	if err := lexicon.CheckPattern("did", r.Did); err != nil {
		return Error(err, PatternErrorCode, "fail did validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("site", r.Site); err != nil {
		return Error(err, PatternErrorCode, "fail site validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("processing", r.Processing); err != nil {
		return Error(err, PatternErrorCode, "fail processing string validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("dataset_parent", r.Parent); err != nil {
		return Error(err, PatternErrorCode, "fail dataset parent validation", "dbs.DatasetRecord.Validate")
	}
	for _, b := range r.Buckets {
		if err := lexicon.CheckPattern("bucket", b); err != nil {
			return Error(err, PatternErrorCode, "fail bucket validation", "dbs.DatasetRecord.Validate")
		}
	}
	for _, f := range r.InputFiles {
		if err := lexicon.CheckPattern("fail", f); err != nil {
			return Error(err, PatternErrorCode, "fail file validation", "dbs.DatasetRecord.Validate")
		}
	}
	for _, f := range r.OutputFiles {
		if err := lexicon.CheckPattern("fail", f); err != nil {
			return Error(err, PatternErrorCode, "fail file validation", "dbs.DatasetRecord.Validate")
		}
	}
	if err := lexicon.CheckPattern("osinfo_name", r.OsInfo.Name); err != nil {
		return Error(err, PatternErrorCode, "fail osinfo.name validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("osinfo_version", r.OsInfo.Version); err != nil {
		return Error(err, PatternErrorCode, "fail osinfo.version validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("osinfo_kernel", r.OsInfo.Kernel); err != nil {
		return Error(err, PatternErrorCode, "fail osinfo.name validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_name", r.Environment.Name); err != nil {
		return Error(err, PatternErrorCode, "fail env.Name validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_version", r.Environment.Version); err != nil {
		return Error(err, PatternErrorCode, "fail env.Version validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_details", r.Environment.Details); err != nil {
		return Error(err, PatternErrorCode, "fail env.Details validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_parent", r.Environment.Parent); err != nil {
		return Error(err, PatternErrorCode, "fail env.Parent validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("script_name", r.Script.Name); err != nil {
		return Error(err, PatternErrorCode, "fail script.Name validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("script_options", r.Script.Options); err != nil {
		return Error(err, PatternErrorCode, "fail script.Options validation", "dbs.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("script_parent", r.Script.Parent); err != nil {
		return Error(err, PatternErrorCode, "fail script.Parent validation", "dbs.DatasetRecord.Validate")
	}
	return nil
}
