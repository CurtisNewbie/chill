alter table chill.command_log drop column build_name;

alter table build_log add column build_start_time TIMESTAMP COMMENT 'build start time',
    add column build_end_time TIMESTAMP COMMENT 'build end time';

update build_log set build_start_time = ctime, build_end_time = utime where build_start_time is null or build_end_time is null;