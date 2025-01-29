--https://stackoverflow.com/questions/18387209/sqlite-syntax-for-creating-table-with-foreign-key

create TABLE processing (
    processing_id INTEGER PRIMARY KEY AUTOINCREMENT,
    processing VARCHAR(255) NOT NULL UNIQUE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
create TABLE parents (
    parent_id INTEGER NOT NULL,
    dataset_id INTEGER NOT NULL,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
create TABLE sites (
    site_id INTEGER PRIMARY KEY AUTOINCREMENT,
    site VARCHAR(255) NOT NULL UNIQUE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
create TABLE buckets (
    bucket_id INTEGER PRIMARY KEY AUTOINCREMENT,
    bucket VARCHAR(255) NOT NULL UNIQUE,
    uuid VARCHAR(255),
    meta_data TEXT,
    dataset_id INTEGER REFERENCES datasets(dataset_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
create TABLE datasets (
    dataset_id INTEGER PRIMARY KEY AUTOINCREMENT,
    did VARCHAR(255) NOT NULL UNIQUE,
    site_id INTEGER REFERENCES sites(site_id) ON UPDATE CASCADE,
    processing_id INTEGER REFERENCES processing(processing_id) ON UPDATE CASCADE,
    os_id INTEGER REFERENCES osinfo(osinfo_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
create TABLE files (
    file_id INTEGER PRIMARY KEY AUTOINCREMENT,
    file VARCHAR(255) NOT NULL UNIQUE,
    checksum VARCHAR(255),
    size INTEGER,
    is_file_valid INTEGER DEFAULT 1,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE osinfo (
    os_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    kernel VARCHAR(255),
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE environments (
    environment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    details TEXT,
    os_id INTEGER,
    parent_environment_id INTEGER,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255),
    FOREIGN KEY (os_id) REFERENCES osinfo(os_id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (parent_environment_id) REFERENCES environments(environment_id) ON DELETE SET NULL ON UPDATE CASCADE
);
CREATE TABLE packages (
    package_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE scripts (
    script_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    order_idx INTEGER,
    options TEXT,
    parent_script_id INTEGER,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255),
    FOREIGN KEY (parent_script_id) REFERENCES scripts(script_id) ON DELETE SET NULL ON UPDATE CASCADE
);
-- Many-to-many relationships

-- dataset may have input and output files, and file can be present in
-- different datasets
CREATE TABLE datasets_files (
    dataset_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    file_type TEXT,
    PRIMARY KEY (dataset_id, file_id, file_type),  -- Prevents duplicate dataset-file-type combinations
    FOREIGN KEY (dataset_id) REFERENCES datasets(dataset_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(file_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- dataset may have many environments, and one environment can be associated
-- with different datasets
CREATE TABLE datasets_environments (
    dataset_id INTEGER NOT NULL,
    environment_id INTEGER NOT NULL,
    PRIMARY KEY (dataset_id, environment_id),  -- Prevents duplicate dataset-environment combinations
    FOREIGN KEY (dataset_id) REFERENCES datasets(dataset_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(environment_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- dataset may have many scripts, and one script can be associated
-- with different datasets
CREATE TABLE datasets_scripts (
    dataset_id INTEGER NOT NULL,
    script_id INTEGER NOT NULL,
    PRIMARY KEY (dataset_id, script_id),  -- Prevents duplicate dataset-script combinations
    FOREIGN KEY (dataset_id) REFERENCES datasets(dataset_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (script_id) REFERENCES scripts(script_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- environment can have multiple python packages and a given package may be
-- presented in different environments
CREATE TABLE environments_packages (
    environment_id INTEGER NOT NULL,
    package_id INTEGER NOT NULL,
    PRIMARY KEY (environment_id, package_id),
    FOREIGN KEY (environment_id) REFERENCES environments(environment_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (package_id) REFERENCES packages(package_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- indexes
CREATE INDEX idx_datasets_did ON datasets(did);
CREATE INDEX idx_files_file ON files(file);
CREATE INDEX idx_scripts_name ON scripts(name);
CREATE INDEX idx_environments_name ON environments(name);
CREATE INDEX idx_packages_name ON packages(name);
CREATE INDEX idx_processing_name ON processing(processing);
CREATE INDEX idx_osinfo_name ON osinfo(name);
CREATE INDEX idx_osinfo_kernel ON osinfo(kernel);
CREATE INDEX idx_osinfo_version ON osinfo(version);

