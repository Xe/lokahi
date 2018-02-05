CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS checks
 ( id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc()
 , created_at TIMESTAMP    NOT NULL DEFAULT NOW()
 , edited_at  TIMESTAMP    NOT NULL DEFAULT NOW()
 , url TEXT UNIQUE NOT NULL
 , webhook_url TEXT NOT NULL
 , playbook_url TEXT NOT NULL
 , every INTEGER DEFAULT 60
 , state TEXT NOT NULL DEFAULT 'init'
 );

CREATE INDEX IF NOT EXISTS checks_every ON checks(every);
