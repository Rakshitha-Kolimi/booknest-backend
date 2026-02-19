CREATE TABLE IF NOT EXISTS authors (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

ALTER TABLE books
ADD COLUMN IF NOT EXISTS author_id UUID;

INSERT INTO authors (name)
SELECT DISTINCT b.author_name
FROM books b
WHERE b.author_name IS NOT NULL
  AND b.author_name <> ''
ON CONFLICT (name) DO NOTHING;

UPDATE books b
SET author_id = a.id
FROM authors a
WHERE b.author_id IS NULL
  AND lower(a.name) = lower(b.author_name);

ALTER TABLE books
ALTER COLUMN author_id SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'fk_books_author'
  ) THEN
    ALTER TABLE books
    ADD CONSTRAINT fk_books_author FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE RESTRICT;
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_books_author_id ON books(author_id);

ALTER TABLE books
DROP COLUMN IF EXISTS author_name;
