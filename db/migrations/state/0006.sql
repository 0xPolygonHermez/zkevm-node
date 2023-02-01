-- +migrate Down
ALTER TABLE state.proof
DROP COLUMN IF EXISTS created_at;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS updated_at;
ALTER TABLE state.proof
ADD COLUMN  generating BOOLEAN DEFAULT FALSE;
UPDATE state.proof
SET generating = TRUE WHERE generating_since IS NOT NULL;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS generating_since;

-- +migrate Up
ALTER TABLE state.proof
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE NOT NULL;
ALTER TABLE state.proof
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL;
ALTER TABLE state.proof
ADD COLUMN generating_since TIMESTAMP WITH TIME ZONE;
UPDATE state.proof
SET generating_since = NOW() WHERE generating IS TRUE;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS generating;
