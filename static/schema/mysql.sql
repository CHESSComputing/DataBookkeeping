CREATE TABLE processing (
    processing_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    processing VARCHAR(255) NOT NULL UNIQUE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

CREATE TABLE parents (
    parent_id INTEGER NOT NULL,
    dataset_id INTEGER NOT NULL,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

CREATE TABLE sites (
    site_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    site VARCHAR(255) NOT NULL UNIQUE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

CREATE TABLE buckets (
    bucket_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    bucket VARCHAR(255) NOT NULL UNIQUE,
    uuid VARCHAR(255),
    meta_data TEXT,
    dataset_id BIGINT REFERENCES datasets(dataset_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

CREATE TABLE datasets (
    dataset_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    did VARCHAR(255) NOT NULL UNIQUE,
    site_id BIGINT REFERENCES sites(site_id) ON UPDATE CASCADE,
    processing_id BIGINT REFERENCES processing(processing_id) ON UPDATE CASCADE,
    os_id BIGINT REFERENCES osinfo(osinfo_id) ON UPDATE CASCADE,
    config_id BIGINT REFERENCES configs(config_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

CREATE TABLE files (
    file_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    file VARCHAR(255) NOT NULL UNIQUE,
    checksum VARCHAR(255),
    size INTEGER,
    is_file_valid INTEGER DEFAULT 1,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

CREATE TABLE osinfo (
    os_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    kernel VARCHAR(255),
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

CREATE TABLE environments (
    environment_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    details TEXT,
    os_id INTEGER,
    parent_environment_id INTEGER,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255),
    FOREIGN KEY (os_id) REFERENCES osinfo(os_id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE packages (
    package_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

-- Scripts
CREATE TABLE scripts (
    script_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,     -- Name of the script
    order_idx INTEGER,
    options TEXT,                   -- Parameters used for the script
    parent_script_id INTEGER,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

-- Configs
CREATE TABLE configs (
    config_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    content TEXT,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

-- Many-to-many relationships

-- dataset may have input and output files, and file can be present in
-- different datasets
CREATE TABLE datasets_files (
    dataset_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    file_type VARCHAR(255),
    PRIMARY KEY (dataset_id, file_id, file_type),  -- Prevents duplicate dataset-file-type combinations
    FOREIGN KEY (dataset_id) REFERENCES datasets(dataset_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(file_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- dataset may have many environments, and one environment can be associated
-- with different datasets
CREATE TABLE datasets_environments (
    dataset_id INTEGER NOT NULL,
    environment_id INTEGER NOT NULL,
    PRIMARY KEY (dataset_id, environment_id),  -- Prevents duplicate dataset-environment combinations
    FOREIGN KEY (dataset_id) REFERENCES datasets(dataset_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(environment_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- dataset may have many scripts, and one script can be associated
-- with different datasets
CREATE TABLE datasets_scripts (
    dataset_id INTEGER NOT NULL,
    script_id INTEGER NOT NULL,
    PRIMARY KEY (dataset_id, script_id),  -- Prevents duplicate dataset-script combinations
    FOREIGN KEY (dataset_id) REFERENCES datasets(dataset_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (script_id) REFERENCES scripts(script_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- dataset may have many configs, and one config can be associated
-- with different datasets
CREATE TABLE datasets_configs (
    dataset_id INTEGER NOT NULL,
    config_id INTEGER NOT NULL,
    PRIMARY KEY (dataset_id, config_id),  -- Prevents duplicate dataset-config combinations
    FOREIGN KEY (dataset_id) REFERENCES datasets(dataset_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (config_id) REFERENCES configs(config_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- environment can have multiple python packages and a given package may be
-- presented in different environments
CREATE TABLE environments_packages (
    environment_id INTEGER NOT NULL,
    package_id INTEGER NOT NULL,
    PRIMARY KEY (environment_id, package_id),
    FOREIGN KEY (environment_id) REFERENCES environments(environment_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (package_id) REFERENCES packages(package_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

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
