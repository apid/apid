PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE APID (
     instance_id text,
     apid_cluster_id text,
     last_snapshot_info text,
     PRIMARY KEY (instance_id)
 );
INSERT INTO "APID" VALUES('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa','1:1:');
COMMIT;
