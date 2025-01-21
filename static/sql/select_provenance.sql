SELECT 
    d.did AS dataset_did, 
    pd.did AS parent_did,
    p.processing,
    o.name AS os_name,
    o.kernel AS os_kernel,
    o.version AS os_version,
    e.environment_id,
    e.name AS environment_name,
    e.version AS environment_version, 
    e.details AS environment_details,
    pe.name AS parent_environment_name,
    eo.name AS env_os_name,
    pk.name AS package_name,
    pk.version AS package_version,
    sc.name AS script_name,
    sc.options AS script_options,
    ps.name AS parent_script_name,
    s.site AS site_name,
    f.file,
    df.file_type,
    b.bucket
FROM datasets d
LEFT JOIN datasets pd ON d.parent_id = pd.dataset_id
LEFT JOIN processing p ON d.processing_id = p.processing_id
LEFT JOIN sites s ON d.site_id = s.site_id
LEFT JOIN datasets_files df ON d.dataset_id = df.dataset_id
LEFT JOIN files f ON df.file_id = f.file_id
LEFT JOIN datasets_environments de ON d.dataset_id = de.dataset_id
LEFT JOIN environments e ON de.environment_id = e.environment_id
LEFT JOIN environments pe ON e.parent_environment_id = pe.environment_id
LEFT JOIN osinfo o ON d.os_id = o.os_id
LEFT JOIN environments eo ON e.os_id = eo.os_id
LEFT JOIN datasets_scripts ds ON d.dataset_id = ds.dataset_id
LEFT JOIN scripts sc ON ds.script_id = sc.script_id
LEFT JOIN scripts ps ON sc.parent_script_id = ps.script_id
LEFT JOIN environments_packages ep ON e.environment_id = ep.environment_id
LEFT JOIN packages pk ON ep.package_id = pk.package_id
LEFT JOIN buckets b ON b.dataset_id = d.dataset_id
