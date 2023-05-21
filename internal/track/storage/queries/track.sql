-- name: CreateTrack :exec
insert into "track" (id, title, duration, files, uploaded_by, uploaded_at)
values ($1, $2, $3, $4, $5, $6);

-- name: FindTracks :many
select * from "track" t
where t.id = ANY (@track_ids::uuid[]);

-- name: PaginateTracks :many
select *
from "track"
where sqlc.narg(after_uploaded_at)::timestamp is null
   or uploaded_at < @after_uploaded_at
   or (sqlc.narg(after_id)::uuid is null or (uploaded_at = @after_uploaded_at and id > @after_id))
order by uploaded_at desc, id
limit @first;
