SELECT DISTINCT
    d.did,
    o.name AS osinfo_name,
    o.version AS osinfo_version,
    o.kernel AS osinfo_kernel,
    o.create_by,
    o.create_at,
    o.modify_by,
    o.modify_at
FROM datasets d
LEFT JOIN osinfo o ON d.os_id = o.os_id
