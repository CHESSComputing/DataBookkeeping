INSERT INTO processing
    (processing_id,processing,
     environment_id,os_id,script_id,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:processing_id,:processing,
     :environment_id,:os_id,:script_id,
     :create_at,:create_by,
     :modify_at,:modify_by)
