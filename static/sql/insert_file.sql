INSERT INTO FILES
    (file_id,logical_file_name,is_file_valid,
     dataset_id,meta_id,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:file_id,:logical_file_name,:is_file_valid,
     :dataset_id,:meta_id,
     :create_at,:create_by,
     :modify_at,:modify_by)
