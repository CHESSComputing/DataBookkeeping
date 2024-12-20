SELECT
    D.did,
    PR.processing,
    S.name AS script_name,
    S.options AS script_options,
    SP.name AS parent_script_name,
    S.create_by,
    S.create_at,
    S.modify_by,
    S.modify_at
FROM scripts S
LEFT JOIN processing PR on PR.script_id=S.script_id
LEFT JOIN datasets D on D.processing_id=PR.processing_id
LEFT JOIN scripts as SP ON S.parent_script_id=SP.script_id
