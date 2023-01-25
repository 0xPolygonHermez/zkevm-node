-- +migrate Down
ALTER TABLE state.proof
DROP COLUMN IF EXISTS created_at;
ALTER TABLE state.proof
DROP COLUMN IF EXISTS updated_at;

-- +migrate Up
ALTER TABLE state.proof
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE NOT NULL;
ALTER TABLE state.proof
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL;
