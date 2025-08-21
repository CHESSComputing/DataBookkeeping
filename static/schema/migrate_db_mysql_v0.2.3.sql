-- 1. Add new table `configs`
CREATE TABLE IF NOT EXISTS configs (
    config_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    content TEXT,
    create_at INTEGER,
    create_by VARCHAR(255),
    modify_at INTEGER,
    modify_by VARCHAR(255)
) ENGINE=InnoDB;

-- 2. Add new column `config_id` to `datasets`
ALTER TABLE datasets
    ADD COLUMN config_id BIGINT NULL;

-- 3. Add foreign key constraint from `datasets.config_id` to `configs.config_id`
ALTER TABLE datasets
    ADD CONSTRAINT fk_datasets_config_id
    FOREIGN KEY (config_id) REFERENCES configs(config_id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;

-- 4. Create join table for many-to-many relationship
CREATE TABLE IF NOT EXISTS datasets_configs (
    dataset_id INTEGER NOT NULL,
    config_id INTEGER NOT NULL,
    PRIMARY KEY (dataset_id, config_id),
    FOREIGN KEY (dataset_id) REFERENCES datasets(dataset_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (config_id) REFERENCES configs(config_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

