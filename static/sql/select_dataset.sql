SELECT DISTINCT
    D.did,
    D.create_by,
    D.create_at,
    D.modify_by,
    D.modify_at
FROM datasets D
{{if .Sites}}
LEFT JOIN sites S on S.site_id=D.site_id
{{end}}
{{if .Processing}}
LEFT JOIN processing PR on PR.processing_id=D.processing_id
{{end}}
{{if .Environments}}
LEFT JOIN datasets_environments DE ON D.dataset_id = DE.dataset_id
LEFT JOIN environments E ON DE.environment_id = E.environment_id
LEFT JOIN environments_packages EP ON e.environment_id = EP.environment_id
LEFT JOIN packages PK ON EP.package_id = PK.package_id
{{end}}
{{if .Osinfo}}
LEFT JOIN osinfo O ON D.os_id = O.os_id
{{end}}
{{if .Scripts}}
LEFT JOIN datasets_scripts DS ON D.dataset_id = DS.dataset_id
LEFT JOIN scripts SC ON DS.script_id = SC.script_id
{{end}}
{{if .Files}}
LEFT JOIN datasets_files DF ON D.dataset_id = DF.dataset_id
LEFT JOIN files F ON F.file_id = DF.file_id
{{end}}
{{if .Buckets}}
LEFT JOIN buckets b ON b.dataset_id = d.dataset_id
{{end}}
