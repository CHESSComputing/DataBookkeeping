SELECT DISTINCT
    d.did,
    f.file AS name,
    f.checksum,
    f.size,
    f.is_file_valid,
    f.create_by,
    f.create_at,
    f.modify_by,
    f.modify_at,
    df.file_type
FROM files f
LEFT JOIN datasets_files df ON df.file_id = f.file_id
LEFT JOIN datasets d on d.dataset_id = df.dataset_id
