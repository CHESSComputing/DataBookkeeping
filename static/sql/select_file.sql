SELECT
    D.did,
    F.file AS name,
    F.checksum,
    F.size,
    F.is_file_valid,
    F.create_by,
    F.create_at,
    F.modify_by,
    F.modify_at,
    DF.file_type
FROM files F
LEFT JOIN datasets_files DF ON DF.file_id = F.file_id
LEFT JOIN datasets D on D.dataset_id = DF.dataset_id
