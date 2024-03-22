INSERT INTO files
    (file_id,file,is_file_valid,
     dataset_id,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:file_id,:file,:is_file_valid,
     :dataset_id,
     :create_at,:create_by,
     :modify_at,:modify_by)
