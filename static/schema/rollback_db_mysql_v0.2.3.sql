-- 1. Drop the many-to-many join table
DROP TABLE IF EXISTS datasets_configs;

-- 2. Remove foreign key constraint and column from datasets
ALTER TABLE datasets
    DROP FOREIGN KEY fk_datasets_config_id;

ALTER TABLE datasets
    DROP COLUMN config_id;

-- 3. Drop the configs table
DROP TABLE IF EXISTS configs;

