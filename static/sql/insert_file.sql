INSERT INTO FILES
    (file_id,logical_file_name,is_file_valid,
     dataset_id,meta_id,
     creation_date,create_by,
     last_modification_date,last_modified_by)
    VALUES
    (:file_id,:logical_file_name,:is_file_valid,
     :dataset_id,:meta_id,
     :creation_date,:create_by,
     :last_modification_date,:last_modified_by)
