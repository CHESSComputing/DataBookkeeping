SELECT DISTINCT
    d.did,
    pr.processing,
    s.name AS script_name,
    s.order_idx AS order_idx,
    s.options AS script_options,
    ps.name AS parent_script_name,
    s.create_by,
    s.create_at,
    s.modify_by,
    s.modify_at
FROM datasets d
LEFT JOIN processing pr on pr.processing_id=d.processing_id
LEFT JOIN datasets_scripts ds ON d.dataset_id = ds.dataset_id
LEFT JOIN scripts s ON ds.script_id = s.script_id
LEFT JOIN scripts ps ON s.parent_script_id = ps.script_id
