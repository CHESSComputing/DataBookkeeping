SELECT
    D.did,
    PR.processing,
    E.name AS environment_name,
    E.version AS environment_version,
    E.details AS environment_details,
    EP.name AS parent_script_name,
    E.create_by,
    E.create_at,
    E.modify_by,
    E.modify_at
FROM environments E
LEFT JOIN processing PR on PR.environment_id=E.environment_id
LEFT JOIN datasets D on D.processing_id=PR.processing_id
LEFT JOIN environments as EP ON E.parent_environment_id=EP.environment_id
