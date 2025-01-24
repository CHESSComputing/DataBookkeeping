SELECT
    D.did,
    PR.processing,
    S.name AS script_name,
    S.order_idx AS order_idx,
    S.options AS script_options,
    PS.name AS parent_script_name,
    S.create_by,
    S.create_at,
    S.modify_by,
    S.modify_at
FROM datasets D
LEFT JOIN processing PR on PR.processing_id=D.processing_id
LEFT JOIN datasets_scripts DS ON D.dataset_id = DS.dataset_id
LEFT JOIN scripts S ON DS.script_id = S.script_id
LEFT JOIN scripts PS ON S.parent_script_id = PS.script_id
