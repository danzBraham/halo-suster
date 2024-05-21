SELECT 
  p.identity_number, p.phone_number, p.name AS patient_name, p.birth_date, p.gender, p.card_image_url,
  m.symptoms, m.medications, m.created_at,
  u.nip, u.name AS user_name, u.id
FROM medical_records m
INNER JOIN patients p ON m.patient_identity_number = p.identity_number
INNER JOIN users u ON m.created_by = u.id;