
-- +migrate Up
ALTER TABLE products ADD COLUMN search_text tsvector;
CREATE INDEX search_text_idx ON products USING GIN (search_text);
UPDATE products SET search_text = (
    setweight(to_tsvector(coalesce(name, '')), 'A') || 
    setweight(to_tsvector(coalesce(description, '')), 'B')
);
-- +migrate Down
ALTER TABLE products DROP COLUMN IF EXISTS search_text;