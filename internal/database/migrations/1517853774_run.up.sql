CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS runs
 ( id         SERIAL
 , uuid       UUID      PRIMARY KEY DEFAULT uuid_generate_v1mc()
 , created_at TIMESTAMP NOT NULL DEFAULT NOW()
 , message    TEXT      NOT NULL DEFAULT ''
 );

CREATE TABLE IF NOT EXISTS run_info
 ( id                                SERIAL
 , created_at                        TIMESTAMP NOT NULL DEFAULT NOW()
 , run_id                            TEXT      NOT NULL
 , check_id                          TEXT      NOT NULL
 , response_time_nanoseconds         BIGINT    NOT NULL
 , webhook_response_time_nanoseconds BIGINT    NOT NULL
 );

CREATE INDEX IF NOT EXISTS run_info_check_id ON run_info(check_id);
