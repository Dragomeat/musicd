-- name: CreatePlayer :exec
insert into "player" (id, host_id, updated_at, created_at)
values ($1, $2, $3, $4);

-- name: FindPlayer :one
select * from "player" where id = $1;

-- name: FindPlayerForUpdate :one
select * from "player" where id = $1 for update;

-- name: UpdatePlayer :exec
update "player" set host_id = $1, current_track = $2, queue = $3, updated_at = $4 where id = $5;
