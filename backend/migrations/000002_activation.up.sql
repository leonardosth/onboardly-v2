-- Add meeting lifecycle status and completion tracking
ALTER TABLE meetings ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'scheduled';
ALTER TABLE meetings ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;

-- Add CHECK constraint for meeting status values
ALTER TABLE meetings ADD CONSTRAINT chk_meeting_status CHECK (status IN ('scheduled', 'completed', 'cancelled'));

-- Add client activation timestamp to projects
ALTER TABLE projects ADD COLUMN IF NOT EXISTS activated_at TIMESTAMPTZ;

-- Indexes for new queries
CREATE INDEX IF NOT EXISTS idx_meetings_analyst_status ON meetings(analyst_id, status);
CREATE INDEX IF NOT EXISTS idx_projects_activated_at ON projects(activated_at);
