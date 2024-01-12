SELECT * FROM files F
JOIN datasets D on D.dataset_id = F.dataset_id
JOIN metadata M on M.meta_id = F.meta_id
