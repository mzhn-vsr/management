CREATE TABLE IF NOT EXISTS feedback (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  question VARCHAR NOT NULL,
  answer VARCHAR NOT NULL,
  is_useful BOOLEAN,
  created_at TIMESTAMP DEFAULT now()
);
