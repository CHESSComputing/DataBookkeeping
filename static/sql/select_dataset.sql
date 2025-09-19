SELECT DISTINCT
    d.did,
    d.create_by,
    d.create_at,
    d.modify_by,
    d.modify_at
FROM datasets d
{{if .Sites}}
LEFT JOIN sites s on s.site_id=d.site_id
{{end}}
{{if .Processing}}
LEFT JOIN processing pr on pr.processing_id=d.processing_id
{{end}}
{{if .Environments}}
LEFT JOIN datasets_environments de ON d.dataset_id = de.dataset_id
LEFT JOIN environments e ON de.environment_id = e.environment_id
LEFT JOIN environments_packages ep ON e.environment_id = ep.environment_id
LEFT JOIN packages pk ON ep.package_id = pk.package_id
{{end}}
{{if .Osinfo}}
LEFT JOIN osinfo o ON d.os_id = o.os_id
{{end}}
{{if .Config}}
LEFT JOIN config c ON d.config_id = c.config_id
{{end}}
{{if .Scripts}}
LEFT JOIN datasets_scripts ds ON d.dataset_id = ds.dataset_id
LEFT JOIN scripts sc ON ds.script_id = sc.script_id
{{end}}
{{if .Files}}
LEFT JOIN datasets_files df ON d.dataset_id = df.dataset_id
LEFT JOIN files f ON f.file_id = df.file_id
{{end}}
{{if .Buckets}}
LEFT JOIN buckets b ON b.dataset_id = d.dataset_id
{{end}}
