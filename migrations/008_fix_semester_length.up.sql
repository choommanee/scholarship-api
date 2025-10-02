-- Fix semester field length in scholarships table
ALTER TABLE scholarships ALTER COLUMN semester TYPE VARCHAR(50);

-- Also fix academic_year if needed
ALTER TABLE scholarships ALTER COLUMN academic_year TYPE VARCHAR(20); 