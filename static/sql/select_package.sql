SELECT DISTINCT
    d.did,
    e.name AS environment_name,
    pk.name AS package_name,
    pk.version AS package_version
FROM datasets D
LEFT JOIN datasets_environments de ON d.dataset_id = de.dataset_id
LEFT JOIN environments e ON de.environment_id = e.environment_id
LEFT JOIN environments eo ON e.os_id = eo.os_id
LEFT JOIN environments_packages ep ON e.environment_id = ep.environment_id
LEFT JOIN packages pk ON ep.package_id = pk.package_id
