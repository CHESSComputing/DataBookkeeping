SELECT DISTINCT
    d.did,
    pr.processing,
    e.name AS environment_name,
    e.version AS environment_version,
    e.details AS environment_details,
    pe.name AS parent_script_name,
    e.create_by,
    e.create_at,
    e.modify_by,
    e.modify_at
FROM datasets d
LEFT JOIN processing pr on pr.processing_id=d.processing_id
LEFT JOIN datasets_environments de ON d.dataset_id = de.dataset_id
LEFT JOIN environments e ON de.environment_id = e.environment_id
LEFT JOIN environments pe ON e.parent_environment_id = pe.environment_id
