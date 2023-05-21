-- name: CreateTrack :exec
insert into "track" (id, title, duration, files, uploaded_by, uploaded_at)
values ($1, $2, $3, $4, $5, $6);

-- name: FindTracks :many
select * from "track" t
where t.id = ANY (@track_ids::uuid[]);
