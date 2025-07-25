SELECT DISTINCT
    D.did,
    C.config,
    C.create_by,
    C.create_at,
    C.modify_by,
    C.modify_at
FROM datasets D
LEFT JOIN configs C ON D.config_id = C.config_id
