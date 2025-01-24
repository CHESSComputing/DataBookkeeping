INSERT INTO buckets
    (bucket_id,bucket,uuid,meta_data,dataset_id,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:bucket_id,:bucket,:uuid,:meta_data,:dataset_id,
     :create_at,:create_by,
     :modify_at,:modify_by)
