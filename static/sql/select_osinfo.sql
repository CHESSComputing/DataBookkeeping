SELECT
    D.did,
    PR.processing,
    O.name AS osinfo_name,
    O.version AS osinfo_version,
    O.kernel AS osinfo_kernel,
    O.create_by,
    O.create_at,
    O.modify_by,
    O.modify_at
FROM osinfo O
LEFT JOIN processing PR on PR.os_id=O.os_id
LEFT JOIN datasets D on D.processing_id=PR.processing_id
