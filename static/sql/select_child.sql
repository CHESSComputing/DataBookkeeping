SELECT DISTINCT
    d.did child_did,
    pds.did,
    d.create_at,
    d.create_by,
    d.modify_at,
    d.modify_by
FROM datasets d
LEFT OUTER JOIN parents dsp ON dsp.dataset_id = d.dataset_id
LEFT OUTER JOIN datasets pds ON pds.dataset_id = dsp.parent_id
