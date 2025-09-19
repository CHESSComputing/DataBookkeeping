CREATE TABLE `processing` (
  `processing_id` int(11) NOT NULL AUTO_INCREMENT,
  `processing` varchar(255) NOT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`processing_id`),
  UNIQUE KEY `processing` (`processing`),
  KEY `idx_processing_name` (`processing`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `parents` (
  `parent_id` int(11) NOT NULL,
  `dataset_id` int(11) NOT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`dataset_id`,`parent_id`),
  KEY `idx_parents_dataset` (`dataset_id`),
  KEY `idx_parents_parent` (`parent_id`),
  CONSTRAINT `fk_parents_dataset` FOREIGN KEY (`dataset_id`) REFERENCES `datasets` (`dataset_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_parents_parent` FOREIGN KEY (`parent_id`) REFERENCES `datasets` (`dataset_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `sites` (
  `site_id` int(11) NOT NULL AUTO_INCREMENT,
  `site` varchar(255) NOT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`site_id`),
  UNIQUE KEY `site` (`site`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `buckets` (
  `bucket_id` int(11) NOT NULL AUTO_INCREMENT,
  `bucket` varchar(255) NOT NULL,
  `uuid` varchar(255) DEFAULT NULL,
  `meta_data` text,
  `dataset_id` int(11) DEFAULT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`bucket_id`),
  UNIQUE KEY `bucket` (`bucket`),
  KEY `fk_buckets_dataset` (`dataset_id`),
  CONSTRAINT `fk_buckets_dataset` FOREIGN KEY (`dataset_id`) REFERENCES `datasets` (`dataset_id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `datasets` (
  `dataset_id` int(11) NOT NULL AUTO_INCREMENT,
  `did` varchar(255) NOT NULL,
  `site_id` int(11) DEFAULT NULL,
  `processing_id` int(11) DEFAULT NULL,
  `os_id` int(11) DEFAULT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  `config_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`dataset_id`),
  UNIQUE KEY `did` (`did`),
  KEY `idx_datasets_did` (`did`),
  KEY `fk_datasets_config` (`config_id`),
  KEY `fk_datasets_processing` (`processing_id`),
  KEY `fk_datasets_os` (`os_id`),
  KEY `fk_datasets_sites` (`site_id`),
  CONSTRAINT `fk_datasets_config` FOREIGN KEY (`config_id`) REFERENCES `configs` (`config_id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_datasets_os` FOREIGN KEY (`os_id`) REFERENCES `osinfo` (`os_id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_datasets_processing` FOREIGN KEY (`processing_id`) REFERENCES `processing` (`processing_id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_datasets_sites` FOREIGN KEY (`site_id`) REFERENCES `sites` (`site_id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `files` (
  `file_id` int(11) NOT NULL AUTO_INCREMENT,
  `file` varchar(255) NOT NULL,
  `checksum` varchar(255) DEFAULT NULL,
  `size` int(11) DEFAULT NULL,
  `is_file_valid` int(11) DEFAULT '1',
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`file_id`),
  UNIQUE KEY `file` (`file`),
  KEY `idx_files_file` (`file`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `osinfo` (
  `os_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `version` varchar(255) DEFAULT NULL,
  `kernel` varchar(255) DEFAULT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`os_id`),
  KEY `idx_osinfo_name` (`name`),
  KEY `idx_osinfo_kernel` (`kernel`),
  KEY `idx_osinfo_version` (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `environments` (
  `environment_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `version` varchar(255) DEFAULT NULL,
  `details` text,
  `os_id` int(11) DEFAULT NULL,
  `parent_environment_id` int(11) DEFAULT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`environment_id`),
  KEY `os_id` (`os_id`),
  KEY `idx_environments_name` (`name`),
  KEY `fk_environments_parent` (`parent_environment_id`),
  CONSTRAINT `environments_ibfk_1` FOREIGN KEY (`os_id`) REFERENCES `osinfo` (`os_id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_environments_parent` FOREIGN KEY (`parent_environment_id`) REFERENCES `environments` (`environment_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `packages` (
  `package_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `version` varchar(255) DEFAULT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`package_id`),
  KEY `idx_packages_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Scripts
CREATE TABLE `scripts` (
  `script_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `order_idx` int(11) DEFAULT NULL,
  `options` text,
  `parent_script_id` int(11) DEFAULT NULL,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`script_id`),
  KEY `idx_scripts_name` (`name`),
  KEY `fk_scripts_parent` (`parent_script_id`),
  CONSTRAINT `fk_scripts_parent` FOREIGN KEY (`parent_script_id`) REFERENCES `scripts` (`script_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Configs
CREATE TABLE `configs` (
  `config_id` int(11) NOT NULL AUTO_INCREMENT,
  `content` text,
  `create_at` int(11) DEFAULT NULL,
  `create_by` varchar(255) DEFAULT NULL,
  `modify_at` int(11) DEFAULT NULL,
  `modify_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`config_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Many-to-many relationships

-- dataset may have input and output files, and file can be present in
-- different datasets
CREATE TABLE `datasets_files` (
  `dataset_id` int(11) NOT NULL,
  `file_id` int(11) NOT NULL,
  `file_type` varchar(255) NOT NULL,
  PRIMARY KEY (`dataset_id`,`file_id`,`file_type`),
  KEY `file_id` (`file_id`),
  CONSTRAINT `datasets_files_ibfk_1` FOREIGN KEY (`dataset_id`) REFERENCES `datasets` (`dataset_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `datasets_files_ibfk_2` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- dataset may have many environments, and one environment can be associated
-- with different datasets
CREATE TABLE `datasets_environments` (
  `dataset_id` int(11) NOT NULL,
  `environment_id` int(11) NOT NULL,
  PRIMARY KEY (`dataset_id`,`environment_id`),
  KEY `environment_id` (`environment_id`),
  CONSTRAINT `datasets_environments_ibfk_1` FOREIGN KEY (`dataset_id`) REFERENCES `datasets` (`dataset_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `datasets_environments_ibfk_2` FOREIGN KEY (`environment_id`) REFERENCES `environments` (`environment_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- dataset may have many scripts, and one script can be associated
-- with different datasets
CREATE TABLE `datasets_scripts` (
  `dataset_id` int(11) NOT NULL,
  `script_id` int(11) NOT NULL,
  PRIMARY KEY (`dataset_id`,`script_id`),
  KEY `script_id` (`script_id`),
  CONSTRAINT `datasets_scripts_ibfk_1` FOREIGN KEY (`dataset_id`) REFERENCES `datasets` (`dataset_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `datasets_scripts_ibfk_2` FOREIGN KEY (`script_id`) REFERENCES `scripts` (`script_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- dataset may have many configs, and one config can be associated
-- with different datasets
CREATE TABLE `datasets_configs` (
  `dataset_id` int(11) NOT NULL,
  `config_id` int(11) NOT NULL,
  PRIMARY KEY (`dataset_id`,`config_id`),
  KEY `datasets_configs_fk_config` (`config_id`),
  CONSTRAINT `datasets_configs_fk_config` FOREIGN KEY (`config_id`) REFERENCES `configs` (`config_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `datasets_configs_fk_dataset` FOREIGN KEY (`dataset_id`) REFERENCES `datasets` (`dataset_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- environment can have multiple python packages and a given package may be
-- presented in different environments
CREATE TABLE `environments_packages` (
  `environment_id` int(11) NOT NULL,
  `package_id` int(11) NOT NULL,
  PRIMARY KEY (`environment_id`,`package_id`),
  KEY `package_id` (`package_id`),
  CONSTRAINT `environments_packages_ibfk_1` FOREIGN KEY (`environment_id`) REFERENCES `environments` (`environment_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `environments_packages_ibfk_2` FOREIGN KEY (`package_id`) REFERENCES `packages` (`package_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
