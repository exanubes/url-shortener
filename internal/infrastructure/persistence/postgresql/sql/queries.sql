-- name: CreateLink :exec
INSERT INTO links (id, url, created_at)
VALUES ($1, $2, $3);

-- name: CreateLinkPolicy :exec
INSERT INTO link_policies (link_id, kind, config)
VALUES ($1, $2, $3);

-- name: LogLinkVisit :exec
INSERT INTO link_visits(link_id, visited_at, ip_address)
VALUES($1, $2, $3);

-- name: GetLink :one
SELECT
    l.id,
    l.url,
    l.created_at,
    l.consumed_at,
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'kind', lp.kind, 
            'config', lp.config
        )
    ) AS policies
FROM links l
INNER JOIN link_policies lp ON l.id = lp.link_id
WHERE l.id = $1
GROUP BY 
    l.id, 
    l.url, 
    l.created_at, 
    l.consumed_at;


-- name: CheckShortCodeExists :one
SELECT EXISTS(SELECT 1 FROM links WHERE id = $1);

-- name: ConsumeSingleUseLink :exec
UPDATE links
SET consumed_at = NOW()
WHERE id = $1 AND consumed_at IS NULL
RETURNING 1;

