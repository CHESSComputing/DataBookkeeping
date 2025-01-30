SELECT DISTINCT
    D.did,
    E.name AS environment_name,
    PK.name AS package_name,
    PK.version AS package_version
FROM datasets D
LEFT JOIN datasets_environments DE ON D.dataset_id = DE.dataset_id
LEFT JOIN environments E ON DE.environment_id = E.environment_id
LEFT JOIN environments EO ON E.os_id = EO.os_id
LEFT JOIN environments_packages EP ON E.environment_id = EP.environment_id
LEFT JOIN packages PK ON EP.package_id = PK.package_id
