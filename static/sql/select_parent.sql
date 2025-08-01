SELECT DISTINCT
    d.did,
    pd.did AS parent_did,
    pd.create_at,
    pd.create_by,
    pd.modify_at,
    pd.modify_by
FROM parents p
LEFT JOIN datasets d ON p.dataset_id = d.dataset_id
LEFT JOIN datasets pd ON p.parent_id = pd.dataset_id
