-- Store the full structured AI result JSON so the frontend can display
-- tasks, decisions, risks, etc. without separate table lookups.
ALTER TABLE meetings ADD COLUMN IF NOT EXISTS ai_result_json TEXT;
