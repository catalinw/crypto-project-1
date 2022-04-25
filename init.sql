create table if not exists challenge
(
    id         serial primary key,
    public_key varchar        not null,
    nonce      varchar unique not null,
    expires_at bigint         not null
);