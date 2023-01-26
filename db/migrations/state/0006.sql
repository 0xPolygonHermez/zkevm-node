-- +migrate Down
ALTER TABLE state.proof
DROP COLUMN IF EXISTS created_at;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS updated_at;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS generating_since;
ALTER TABLE state.proof
ADD COLUMN  generating BOOLEAN DEFAULT FALSE;

-- +migrate Up
ALTER TABLE state.proof
DROP COLUMN IF EXISTS generating;
ALTER TABLE state.proof
ADD COLUMN generating_since TIMESTAMP WITH TIME ZONE;
ALTER TABLE state.proof
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE NOT NULL;
ALTER TABLE state.proof
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL;
