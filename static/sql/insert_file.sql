INSERT INTO files
    (file_id,file,checksum,size,is_file_valid,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:file_id,:file,:checksum,:size,:is_file_valid,
     :create_at,:create_by,
     :modify_at,:modify_by)
