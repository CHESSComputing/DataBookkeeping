SELECT DISTINCT
    D.did,
    PR.processing,
    E.name AS environment_name,
    E.version AS environment_version,
    E.details AS environment_details,
    PE.name AS parent_script_name,
    E.create_by,
    E.create_at,
    E.modify_by,
    E.modify_at
FROM datasets D
LEFT JOIN processing PR on PR.processing_id=D.processing_id
LEFT JOIN datasets_environments DE ON D.dataset_id = DE.dataset_id
LEFT JOIN environments E ON DE.environment_id = E.environment_id
LEFT JOIN environments PE ON E.parent_environment_id = PE.environment_id
