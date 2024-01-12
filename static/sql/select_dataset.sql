SELECT
    D.dataset,
    D.meta_id,
    S.site,
    PR.processing,
    P.parent,
    D.create_by,
    D.create_at,
    D.modify_by,
    D.modify_at
FROM datasets D
JOIN sites S on S.site_id=D.site_id
JOIN processing PR on PR.processing_id=D.processing_id
LEFT OUTER JOIN parents P on P.parent_id=D.parent_id
