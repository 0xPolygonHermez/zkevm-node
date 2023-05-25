CREATE TYPE level_t AS ENUM ('emerg', 'alert', 'crit', 'err', 'warning', 'notice', 'info', 'debug');

CREATE TABLE public.event (
   id BIGSERIAL PRIMARY KEY,
   received_at timestamp WITH TIME ZONE default CURRENT_TIMESTAMP,
   ip_address inet,
   source varchar(32) not null,
   component varchar(32),
   level level_t not null,
   event_id varchar(32) not null,
   description text,
   data bytea,
   json jsonb
);
