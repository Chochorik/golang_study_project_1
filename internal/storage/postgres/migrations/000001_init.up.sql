CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_done BOOLEAN NOT NULL DEFAULT false,
    done_at TIMESTAMPTZ,
    notify_status TEXT NOT NULL DEFAULT 'pending'
                   CHECK ( notify_status IN ('pending', 'queued', 'sent', 'failed') ),
    notify_attempts INT NOT NULL DEFAULT 0,
    notified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
)