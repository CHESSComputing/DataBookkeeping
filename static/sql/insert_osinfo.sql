INSERT INTO osinfo
    (osinfo_id,os_name,os_version,kernel_number,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:osinfo_id,:os_name,:os_version,:kernel_number,
     :create_at,:create_by,
     :modify_at,:modify_by)
