CREATE TABLE view_crossfire_chain_stats
(
    metric VARCHAR  NOT NULL,
    value BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (metric)
)

