INSERT INTO scripts
    (script_id,name,order_idx,options,parent_script_id,
     create_at,create_by,
     modify_at,modify_by)
    VALUES
    (:script_id,:name,:order_idx,:options,:parent_script_id,
     :create_at,:create_by,
     :modify_at,:modify_by)
