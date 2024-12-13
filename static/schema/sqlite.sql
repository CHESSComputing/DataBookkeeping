--https://stackoverflow.com/questions/18387209/sqlite-syntax-for-creating-table-with-foreign-key

create TABLE processing (
    processing_id INTEGER PRIMARY KEY AUTOINCREMENT,
    processing VARCHAR(255) NOT NULL UNIQUE,
    environment_id INTEGER REFERENCES python_environments(environment_id) ON UPDATE CASCADE,
    os_id INTEGER REFERENCES os_environments(os_id) ON UPDATE CASCADE,
    script_id INTEGER REFERENCES scripts(script_id) ON UPDATE CASCADE,
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
    processing_id INTEGER REFERENCES processingS(processing_id) ON UPDATE CASCADE,
    parent_id INTEGER REFERENCES parents(parent_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
create TABLE files (
    file_id INTEGER PRIMARY KEY AUTOINCREMENT,
    file VARCHAR(255) NOT NULL UNIQUE,
    is_file_valid INTEGER DEFAULT 1,
    dataset_id INTEGER REFERENCES datasets(dataset_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE python_environments (
    environment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    details TEXT,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE os_environments (
    os_id INTEGER PRIMARY KEY AUTOINCREMENT,
    os_name VARCHAR(255) NOT NULL,
    os_version VARCHAR(255),
    release_number VARCHAR(255),
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE scripts (
    script_id INTEGER PRIMARY KEY AUTOINCREMENT,
    script_name VARCHAR(255) NOT NULL,
    parameters TEXT,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
