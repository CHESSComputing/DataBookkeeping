SELECT 
    d.did AS dataset_did, 
    pd.did AS parent_did
FROM parents p
LEFT JOIN datasets d ON p.dataset_id = d.dataset_id
LEFT JOIN datasets pd ON p.parent_id = pd.dataset_id
--SELECT
--    D.did,
--    PDS.did PARENT_DID,
--FROM datasets D
--LEFT OUTER JOIN parents DSP ON DSP.DATASET_ID = D.DATASET_ID
--LEFT OUTER JOIN datasets PDS ON PDS.DATASET_ID = DSP.PARENT_ID
