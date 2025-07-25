SELECT DISTINCT
    D.did,
    C.content,
    C.create_by,
    C.create_at,
    C.modify_by,
    C.modify_at
FROM datasets D
LEFT JOIN datasets_configs DC ON D.dataset_id = DC.dataset_id
LEFT JOIN configs C ON DC.config_id = C.config_id
