-- +goose Up
CREATE TABLE IF NOT EXISTS companies (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    tax_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employees (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    employment_type TEXT NOT NULL,
    salary_type TEXT NOT NULL,
    base_salary NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_company FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payroll_runs (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL,
    run_date DATE NOT NULL,
    total_employees INT NOT NULL,
    total_pay NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_payroll_company FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payroll_details (
    id BIGSERIAL PRIMARY KEY,
    payroll_run_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,
    gross_pay NUMERIC(12,2) NOT NULL,
    tax_deduction NUMERIC(12,2) DEFAULT 0,
    net_pay NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_payroll_run FOREIGN KEY (payroll_run_id)
        REFERENCES payroll_runs(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_employee FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_employees_company_id ON employees(company_id);
CREATE INDEX idx_payroll_runs_company_id ON payroll_runs(company_id);
CREATE INDEX idx_payroll_details_payroll_run_id ON payroll_details(payroll_run_id);
CREATE INDEX idx_payroll_details_employee_id ON payroll_details(employee_id);

CREATE UNIQUE INDEX idx_payroll_details_run_employee ON payroll_details(payroll_run_id, employee_id);

CREATE INDEX idx_payroll_runs_company_run_date ON payroll_runs(company_id, run_date);

-- +goose Down
DROP TABLE IF EXISTS payroll_details;
DROP TABLE IF EXISTS payroll_runs;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS companies;
