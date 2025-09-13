SELECT DISTINCT
    D.did AS dataset_did, 
    P.processing,
    O.name AS os_name,
    O.kernel AS os_kernel,
    O.version AS os_version,
    E.environment_id,
    E.name AS environment_name,
    E.version AS environment_version, 
    E.details AS environment_details,
    PE.name AS parent_environment_name,
    EO.name AS env_os_name,
    PK.name AS package_name,
    PK.version AS package_version,
    SC.script_id,
    SC.name AS script_name,
    SC.order_idx AS order_idx,
    SC.options AS script_options,
    PS.name AS parent_script_name,
    S.site AS site_name,
    C.content AS config_content,
    B.bucket AS bucket_name,
    B.uuid,
    B.meta_data
FROM datasets D
LEFT JOIN processing P ON D.processing_id = P.processing_id
LEFT JOIN sites S ON D.site_id = S.site_id

-- configs
LEFT JOIN configs C ON D.config_id = C.config_id
LEFT JOIN datasets_configs DC ON D.dataset_id = DC.dataset_id

-- environments
LEFT JOIN datasets_environments DE ON D.dataset_id = DE.dataset_id
LEFT JOIN environments E ON DE.environment_id = E.environment_id
LEFT JOIN environments PE ON E.parent_environment_id = PE.environment_id
LEFT JOIN osinfo EO ON E.os_id = EO.os_id

-- OS
LEFT JOIN osinfo O ON D.os_id = O.os_id

-- scripts
LEFT JOIN datasets_scripts DS ON D.dataset_id = DS.dataset_id
LEFT JOIN scripts SC ON DS.script_id = SC.script_id
LEFT JOIN scripts PS ON SC.parent_script_id = PS.script_id

-- packages
LEFT JOIN environments_packages EP ON E.environment_id = EP.environment_id
LEFT JOIN packages PK ON EP.package_id = PK.package_id

-- buckets
LEFT JOIN buckets B ON B.dataset_id = D.dataset_id
