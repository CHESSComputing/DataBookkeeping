SELECT
    D.did,
    S.site,
    PR.processing,
    E.name AS environment_name,
    E.version AS environment_version,
    E.details AS environment_details,
    O.name AS osinfo_name,
    O.version AS osinfo_version,
    O.kernel AS osinfo_kernel,
    SC.name AS script_name,
    SC.options AS script_options,
    D.create_by,
    D.create_at,
    D.modify_by,
    D.modify_at
FROM datasets D
LEFT JOIN sites S on S.site_id=D.site_id
LEFT JOIN processing PR on PR.processing_id=D.processing_id
LEFT JOIN datasets_environments DE ON D.dataset_id = DE.dataset_id
LEFT JOIN environments E ON DE.environment_id = E.environment_id
LEFT JOIN environments PE ON E.parent_environment_id = PE.environment_id
LEFT JOIN osinfo O ON E.os_id = O.os_id
LEFT JOIN osinfo EO ON E.os_id = EO.os_id
LEFT JOIN datasets_scripts DS ON D.dataset_id = DS.dataset_id
LEFT JOIN scripts SC ON DS.script_id = SC.script_id
LEFT JOIN scripts PS ON SC.parent_script_id = PS.script_id
