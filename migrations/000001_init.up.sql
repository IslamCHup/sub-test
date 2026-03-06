CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,

    service_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),

    user_id UUID NOT NULL,

    start_date DATE NOT NULL,
    end_date DATE,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_subscriptions_user_id
ON subscriptions(user_id);

CREATE INDEX idx_subscriptions_service_name
ON subscriptions(service_name);

CREATE INDEX idx_subscriptions_user_service
ON subscriptions(user_id, service_name);