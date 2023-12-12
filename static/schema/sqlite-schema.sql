--------------------------------------------------------
--  DDL for Table PROCESSING
--------------------------------------------------------

CREATE TABLE "PROCESSING" (
    "PROCESSING_ID" INTEGER,
    "PROCESSING" VARCHAR2(700) NOT NULL UNIQUE,
    "CREATION_DATE" INTEGER,
    "CREATE_BY" VARCHAR2(500),
    "LAST_MODIFICATION_DATE" INTEGER,
    "LAST_MODIFIED_BY" VARCHAR2(500)
);
--------------------------------------------------------
--  DDL for Table PARENTS
--------------------------------------------------------

CREATE TABLE "PARENTS" (
    "PARENT_ID" INTEGER,
    "PARENT" VARCHAR2(700) NOT NULL UNIQUE,
    "CREATION_DATE" INTEGER,
    "CREATE_BY" VARCHAR2(500),
    "LAST_MODIFICATION_DATE" INTEGER,
    "LAST_MODIFIED_BY" VARCHAR2(500)
);
--------------------------------------------------------
--  DDL for Table SITES
--------------------------------------------------------

CREATE TABLE "SITES" (
    "SITE_ID" INTEGER,
    "SITE" VARCHAR2(700) NOT NULL UNIQUE,
    "CREATION_DATE" INTEGER,
    "CREATE_BY" VARCHAR2(500),
    "LAST_MODIFICATION_DATE" INTEGER,
    "LAST_MODIFIED_BY" VARCHAR2(500)
);
--------------------------------------------------------
--  DDL for Table BUCKETS
--------------------------------------------------------

CREATE TABLE "BUCKETS" (
    "BUCKET_ID" INTEGER,
    "BUCKET" VARCHAR2(700) NOT NULL UNIQUE,
    "META_ID" VARCHAR2(700),
    "DATASET_ID" VARCHAR2(700) NOT NULL UNIQUE,
    "CREATION_DATE" INTEGER,
    "CREATE_BY" VARCHAR2(500),
    "LAST_MODIFICATION_DATE" INTEGER,
    "LAST_MODIFIED_BY" VARCHAR2(500)
);
--------------------------------------------------------
--  DDL for Table DATASETS
--------------------------------------------------------

CREATE TABLE "DATASETS" (
    "DATASET_ID" INTEGER,
    "DATASET" VARCHAR2(700) NOT NULL UNIQUE,
    "META_ID" VARCHAR2(700),
    "SITE_ID" INTEGER,
    "PROCESSING_ID" INTEGER,
    "PARENT_ID" INTEGER,
    "CREATION_DATE" INTEGER,
    "CREATE_BY" VARCHAR2(500),
    "LAST_MODIFICATION_DATE" INTEGER,
    "LAST_MODIFIED_BY" VARCHAR2(500)
);
--------------------------------------------------------
--  DDL for Table FILES
--------------------------------------------------------

CREATE TABLE "FILES" (
    "FILE_ID" INTEGER,
    "LOGICAL_FILE_NAME" VARCHAR2(700) NOT NULL UNIQUE,
    "IS_FILE_VALID" INTEGER DEFAULT 1,
    "DATASET_ID" INTEGER,
    "META_ID" VARCHAR2(700),
    "CREATION_DATE" INTEGER,
    "CREATE_BY" VARCHAR2(500),
    "LAST_MODIFICATION_DATE" INTEGER,
    "LAST_MODIFIED_BY" VARCHAR2(500)
);
