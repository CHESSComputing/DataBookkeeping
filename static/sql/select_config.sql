SELECT DISTINCT
    d.did,
    c.content,
    c.create_by,
    c.create_at,
    c.modify_by,
    c.modify_at
FROM datasets d
LEFT JOIN datasets_configs dc ON d.dataset_id = dc.dataset_id
LEFT JOIN configs c ON dc.config_id = c.config_id
