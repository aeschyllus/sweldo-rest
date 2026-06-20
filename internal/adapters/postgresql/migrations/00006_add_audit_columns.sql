-- +goose Up
-- +goose StatementBegin

-- companies (already has updated_at)
ALTER TABLE companies ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
ALTER TABLE companies ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);

-- employees (already has updated_at)
ALTER TABLE employees ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);

-- users (missing updated_at entirely)
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE users ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);

-- payroll_runs (missing updated_at entirely)
ALTER TABLE payroll_runs ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE payroll_runs ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
ALTER TABLE payroll_runs ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);

-- payroll_details (missing updated_at entirely)
ALTER TABLE payroll_details ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE payroll_details ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
ALTER TABLE payroll_details ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);

-- deductions (missing updated_at entirely)
ALTER TABLE deductions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE deductions ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
ALTER TABLE deductions ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE deductions DROP COLUMN IF EXISTS created_by;
ALTER TABLE deductions DROP COLUMN IF EXISTS updated_by;
ALTER TABLE deductions DROP COLUMN IF EXISTS updated_at;

ALTER TABLE payroll_details DROP COLUMN IF EXISTS created_by;
ALTER TABLE payroll_details DROP COLUMN IF EXISTS updated_by;
ALTER TABLE payroll_details DROP COLUMN IF EXISTS updated_at;

ALTER TABLE payroll_runs DROP COLUMN IF EXISTS created_by;
ALTER TABLE payroll_runs DROP COLUMN IF EXISTS updated_by;
ALTER TABLE payroll_runs DROP COLUMN IF EXISTS updated_at;

ALTER TABLE users DROP COLUMN IF EXISTS created_by;
ALTER TABLE users DROP COLUMN IF EXISTS updated_by;
ALTER TABLE users DROP COLUMN IF EXISTS updated_at;

ALTER TABLE employees DROP COLUMN IF EXISTS created_by;
ALTER TABLE employees DROP COLUMN IF EXISTS updated_by;

ALTER TABLE companies DROP COLUMN IF EXISTS created_by;
ALTER TABLE companies DROP COLUMN IF EXISTS updated_by;

-- +goose StatementEnd
