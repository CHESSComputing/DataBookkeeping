INSERT INTO environments
    (environment_id,name,version,details,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:environment_id,:name,:version,:details,
     :create_at,:create_by,
     :modify_at,:modify_by)
