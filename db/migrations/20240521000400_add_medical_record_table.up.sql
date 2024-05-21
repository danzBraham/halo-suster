CREATE TABLE IF NOT EXISTS medical_records (
  id VARCHAR(26) NOT NULL PRIMARY KEY,
  symptoms TEXT NOT NULL,
  medications TEXT NOT NULL,
  patient_identity_number VARCHAR(16) NOT NULL,
  created_by VARCHAR(26) NOT NULL,
  is_deleted BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (patient_identity_number) REFERENCES patients(identity_number),
  FOREIGN KEY (created_by) REFERENCES users(id)
);
