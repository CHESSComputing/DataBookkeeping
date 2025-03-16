SELECT DISTINCT
    D.did CHILD_DID,
    PDS.did,
    D.create_at,
    D.create_by,
    D.modify_at,
    D.modify_by
FROM datasets D
LEFT OUTER JOIN parents DSP ON DSP.DATASET_ID = D.DATASET_ID
LEFT OUTER JOIN datasets PDS ON PDS.DATASET_ID = DSP.PARENT_ID
