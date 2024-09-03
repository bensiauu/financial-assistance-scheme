ALTER TABLE applicants
ADD COLUMN marital_status VARCHAR(50) NOT NULL DEFAULT 'single',
ADD COLUMN disability_status VARCHAR(50) NOT NULL DEFAULT 'none',
ADD COLUMN number_of_children INTEGER NOT NULL DEFAULT 0;
