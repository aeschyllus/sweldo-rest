ALTER TABLE payroll_details ADD COLUMN IF NOT EXISTS hourly_rate NUMERIC(12,2) DEFAULT 0;
ALTER TABLE payroll_details ADD COLUMN IF NOT EXISTS hours_worked NUMERIC(12,2) DEFAULT 0;

CREATE TABLE IF NOT EXISTS deductions (
    id BIGSERIAL PRIMARY KEY,
    payroll_detail_id BIGINT NOT NULL REFERENCES payroll_details(id) ON DELETE CASCADE,
    deduction_type TEXT NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
