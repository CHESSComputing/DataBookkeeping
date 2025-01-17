SELECT 
    d.did AS dataset_did, 
    d.parent_id AS parent_dataset_id,

    -- Processing
    p.processing,

    -- OS Info
    o.name AS os_name,
    o.kernel AS os_kernel,
    o.version AS os_version,

    -- Environment
    e.name AS environment_name,
    e.version AS environment_version, 
    e.details AS environment_details,
    e.parent_environment_id,

    -- Environment OS
    eo.name AS env_os_name,

    -- Packages
    pk.name AS package_name,
    pk.version AS package_version,

    -- Script
    sc.name AS script_name,
    sc.options AS script_options,
    sc.parent_script_id,

    -- Site
    s.site AS site_name,

    -- Files
    f.file,
    df.file_type,

    -- Buckets
    b.bucket

FROM datasets d

LEFT JOIN processing p ON d.processing_id = p.processing_id
LEFT JOIN sites s ON d.site_id = s.site_id
LEFT JOIN dataset_files df ON d.dataset_id = df.dataset_id
LEFT JOIN files f ON df.file_id = f.file_id
LEFT JOIN dataset_environments de ON d.dataset_id = de.dataset_id
LEFT JOIN environments e ON de.environment_id = e.environment_id
LEFT JOIN osinfo o ON e.os_id = o.os_id
LEFT JOIN osinfo eo ON e.os_id = eo.os_id  -- OS for Environment
LEFT JOIN scripts sc ON p.script_id = sc.script_id
LEFT JOIN environment_packages ep ON e.environment_id = ep.environment_id
LEFT JOIN packages pk ON ep.package_id = pk.package_id
LEFT JOIN buckets b ON b.dataset_id = d.dataset_id
