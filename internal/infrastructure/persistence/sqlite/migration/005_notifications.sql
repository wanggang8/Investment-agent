CREATE TABLE IF NOT EXISTS notifications (
  notification_id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  severity TEXT NOT NULL CHECK (severity IN ('info', 'warning', 'critical')),
  title TEXT NOT NULL,
  message TEXT NOT NULL,
  source_type TEXT,
  source_id TEXT,
  read_at DATETIME,
  created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_read_at ON notifications(read_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_notifications_active_source ON notifications(type, source_type, source_id) WHERE read_at IS NULL AND source_type IS NOT NULL AND source_id IS NOT NULL;
