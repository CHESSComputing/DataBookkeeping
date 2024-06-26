CREATE TABLE processing (
    processing_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    processing VARCHAR(255) NOT NULL UNIQUE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE parents (
    parent_id INTEGER NOT NULL,
    dataset_id INTEGER NOT NULL,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE sites (
    site_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    site VARCHAR(255) NOT NULL UNIQUE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE buckets (
    bucket_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    bucket VARCHAR(255) NOT NULL UNIQUE,
    dataset_id BIGINT REFERENCES datasets(dataset_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE datasets (
    dataset_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    did VARCHAR(255) NOT NULL UNIQUE,
    site_id BIGINT REFERENCES site(site_id) ON UPDATE CASCADE,
    processing_id BIGINT REFERENCES processingS(processing_id) ON UPDATE CASCADE,
    parent_id BIGINT REFERENCES parents(parent_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
CREATE TABLE files (
    file_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    file VARCHAR(255) NOT NULL UNIQUE,
    is_file_valid INTEGER DEFAULT 1,
    dataset_id BIGINT REFERENCES datasets(dataset_id) ON UPDATE CASCADE,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
);
