INSERT INTO environments
    (environment_id,name,version,details,os_id,parent_environment_id,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:environment_id,:name,:version,:details,:os_id,:parent_environment_id,
     :create_at,:create_by,
     :modify_at,:modify_by)
