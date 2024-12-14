INSERT INTO scripts
    (script_id,script_name,parameters,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:script_id,:script_name,:parameters,
     :create_at,:create_by,
     :modify_at,:modify_by)
