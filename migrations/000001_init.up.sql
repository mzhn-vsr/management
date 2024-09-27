CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS faq (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  question TEXT NOT NULL,
  answer TEXT NOT NULL,
  classifier1 VARCHAR,
  classifier2 VARCHAR,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP
);


-- CREATE TABLE IF NOT EXISTS docs (
--   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--   name VARCHAR NOT NULL
--   created_at TIMESTAMP NOT NULL DEFAULT now(),
--   updated_at TIMESTAMP
-- );

-- CREATE TABLE IF NOT EXISTS docs_text (
--   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--   document_id UUID REFERENCES docs(id),
--   seq VARCHAR,
--   text TEXT NOT NULL,
--   created_at TIMESTAMP NOT NULL DEFAULT now(),
--   updated_at TIMESTAMP
-- );
