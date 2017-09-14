PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE APID (
     instance_id text,
     apid_cluster_id text,
     last_snapshot_info text,
     PRIMARY KEY (instance_id)
 );
INSERT INTO "APID" VALUES('58c270cb-2a88-4be5-bc41-db53102dcdd5','950b30f1-8c41-4bf5-94a3-f10c104ff5d4','1:1:');
COMMIT;
