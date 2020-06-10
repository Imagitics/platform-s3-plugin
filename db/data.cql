CREATE KEYSPACE platform_s3_db with replication = {'class':'SimpleStrategy', 'replication_factor':1};
use platform_s3_db;
CREATE TABLE platform_s3_db.aws_metadata (tenant_id text PRIMARY KEY, access_key text, preferred_region text, secret_key text);
insert into platform_s3_db.aws_metadata(tenant_id,access_key,preferred_region,secret_key) VALUES ('TEST01', 'AKIA5Z5AMYYVSQDUDBKN','ap-south-1','NEk6/KhDQ9Os70wu7bI5J2vOYUJZOsPalEl6YHBF') ;

