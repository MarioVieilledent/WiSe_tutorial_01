--------------------------------------------------------
--  DDL for Table OSM_WAY_NODE
--------------------------------------------------------

  CREATE TABLE "OSM_NODE" 
   (	"NODE_ID" NUMBER, 
	"LON" NUMBER, 
	"LAT" NUMBER
   )  ;
--------------------------------------------------------
--  DDL for Table OSM_WAY_NODE
--------------------------------------------------------

  CREATE TABLE "OSM_WAY_NODE" 
   (	"WAY_ID" NUMBER, 
	"NODE_ID" NUMBER, 
	"SEQ_NR" NUMBER
   ) ;
--------------------------------------------------------
--  DDL for Table OSM_WAY_TAG
--------------------------------------------------------

  CREATE TABLE "OSM_WAY_TAG" 
   (	"WAY_ID" NUMBER, 
	"KEY" VARCHAR2(100 BYTE), 
	"VALUE" VARCHAR2(250 BYTE)
   ) ;
