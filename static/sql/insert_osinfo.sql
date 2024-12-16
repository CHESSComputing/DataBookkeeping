INSERT INTO osinfo
    (os_id,name,version,kernel,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:os_id,:name,:version,:kernel,
     :create_at,:create_by,
     :modify_at,:modify_by)
