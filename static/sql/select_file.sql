SELECT
    D.did,
    F.file,
    F.is_file_valid,
    F.create_by,
    F.create_at,
    F.modify_by,
    F.modify_at
FROM files F
JOIN datasets D on D.dataset_id = F.dataset_id
LEFT JOIN datasets_files DF ON D.dataset_id = DF.dataset_id
