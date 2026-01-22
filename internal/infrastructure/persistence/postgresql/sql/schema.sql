-- Links table: Core URL mapping
CREATE TABLE IF NOT EXISTS links (
    id VARCHAR(11) PRIMARY KEY,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ
);

-- Link policies table: Expiration policy configuration
CREATE TABLE IF NOT EXISTS link_policies (
    id BIGSERIAL PRIMARY KEY,
    link_id VARCHAR(11) NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    kind TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()

    CONSTRAINT kind CHECK(kind IN ('max_age', 'single_use'))
);

-- Link analytics table: Link tracking data
CREATE TABLE IF NOT EXISTS link_visits (
    id BIGSERIAL PRIMARY KEY,
    link_id VARCHAR(11) NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    visited_at TIMESTAMPTZ NOT NULL,
    ip_address INET
);

CREATE UNIQUE INDEX unique_link_policy_kind ON link_policies(link_id, kind);
