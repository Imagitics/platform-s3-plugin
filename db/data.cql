use platform_s3_db;
CREATE TABLE aws_metadata (tenant_id text PRIMARY KEY,   access_key text,  preferred_region text,   secret_key text);
insert into aws_metadata(tenant_id,access_key,preferred_region,secret_key) VALUES ('TEST01', 'AKIAIHJZ6JMA6LJ4KQBQ','ap-south-1','pv2/HpVjBFEi4fUGMb6oXuaA6ULh4') ;
