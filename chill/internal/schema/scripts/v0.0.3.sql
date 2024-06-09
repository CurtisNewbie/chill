alter table build_log add column tag VARCHAR(45) NOT NULL DEFAULT '' COMMENT 'tag';
alter table build_info add column tag VARCHAR(45) NOT NULL DEFAULT '' COMMENT 'tag';
