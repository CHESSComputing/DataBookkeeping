SELECT DISTINCT
    d.did AS dataset_did, 
    pd.did AS parent_did
FROM parents p
LEFT JOIN datasets d ON p.dataset_id = d.dataset_id
LEFT JOIN datasets pd ON p.parent_id = pd.dataset_id
--SELECT
--    d.did,
--    pds.did parent_did,
--FROM datasets d
--LEFT OUTER JOIN parents dsp ON dsp.dataset_id = d.dataset_id
--LEFT OUTER JOIN datasets pds ON pds.dataset_id = dsp.parent_id
