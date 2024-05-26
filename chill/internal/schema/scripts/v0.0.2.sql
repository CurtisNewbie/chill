alter table build_info add column commit_id VARCHAR(45) NOT NULL DEFAULT '' COMMENT 'commit id';
alter table build_log add column commit_id VARCHAR(45) NOT NULL DEFAULT '' COMMENT 'commit id';