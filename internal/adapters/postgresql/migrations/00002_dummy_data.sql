-- +goose Up
INSERT INTO companies (name, tax_id)
VALUES 
    ('TechSolutions Inc.', '123-456-789-000'),
    ('Green Garden Foods', '987-654-321-001');

INSERT INTO employees (company_id, first_name, last_name, employment_type, salary_type, base_salary)
VALUES 
    ((SELECT id FROM companies WHERE name = 'TechSolutions Inc.'), 'Juan', 'Dela Cruz', 'FULL_TIME', 'MONTHLY', 45000.00),
    ((SELECT id FROM companies WHERE name = 'TechSolutions Inc.'), 'Maria', 'Santos', 'PART_TIME', 'HOURLY', 250.00),
    ((SELECT id FROM companies WHERE name = 'Green Garden Foods'), 'Ricardo', 'Dalisay', 'FULL_TIME', 'MONTHLY', 35000.00),
    ((SELECT id FROM companies WHERE name = 'Green Garden Foods'), 'Elena', 'Adarna', 'CONTRACTOR', 'PROJECT_BASED', 60000.00);

-- +goose Down
DELETE FROM companies WHERE name IN ('TechSolutions Inc.', 'Green Garden Foods');