CREATE INDEX idx_tasks_due ON tasks (done_at)
    WHERE is_done = false AND notify_status = 'pending';