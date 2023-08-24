create table if not exists calendar_event(
  id bigserial not null primary key,
  user_id bigint not null,
  title varchar(255) not null,
  event_time TIMESTAMP DEFAULT NULL,
  duration bigint not null,
  description varchar(255) not null
);