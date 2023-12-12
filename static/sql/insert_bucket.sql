INSERT INTO BUCKETS
    (bucket_id,bucket,meta_id,dataset_id,
     creation_date,create_by,
     last_modification_date,last_modified_by)
    VALUES
    (:bucket_id,:bucket,:meta_id,:dataset_id,
     :creation_date,:create_by,
     :last_modification_date,:last_modified_by)
