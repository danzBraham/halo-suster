ALTER TABLE medical_records
  DROP CONSTRAINT IF EXISTS fk_created_by,
  DROP CONSTRAINT IF EXISTS fk_patient_identity_number;

DROP TABLE IF EXISTS medical_records;
