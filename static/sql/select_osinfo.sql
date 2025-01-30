SELECT DISTINCT
    D.did,
    O.name AS osinfo_name,
    O.version AS osinfo_version,
    O.kernel AS osinfo_kernel,
    O.create_by,
    O.create_at,
    O.modify_by,
    O.modify_at
FROM datasets D
LEFT JOIN osinfo O ON D.os_id = O.os_id
