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
LEFT JOIN environments E on E.environment_id=PR.environment_id
LEFT JOIN osinfo O on O.os_id=PR.os_id
LEFT JOIN scripts SC on SC.script_id=PR.script_id
