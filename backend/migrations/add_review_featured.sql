-- Add is_featured column to reviews table
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS is_featured BOOLEAN DEFAULT false;

-- Create index on is_featured for faster queries
CREATE INDEX IF NOT EXISTS idx_reviews_is_featured ON reviews(is_featured);
