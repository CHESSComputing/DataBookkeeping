SELECT
    D.did,
    S.site,
    PR.processing,
    D.create_by,
    D.create_at,
    D.modify_by,
    D.modify_at
FROM datasets D
LEFT JOIN sites S on S.site_id=D.site_id
LEFT JOIN processing PR on PR.processing_id=D.processing_id
