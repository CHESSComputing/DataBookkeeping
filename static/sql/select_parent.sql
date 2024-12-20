SELECT
    D.did,
    PDS.did PARENT_DID,
    PDS.create_at,
    PDS.create_by,
    PDS.modify_at,
    PDS.modify_by
FROM datasets D
LEFT OUTER JOIN parents DSP ON DSP.DATASET_ID = D.DATASET_ID
LEFT OUTER JOIN datasets PDS ON PDS.DATASET_ID = DSP.PARENT_ID
