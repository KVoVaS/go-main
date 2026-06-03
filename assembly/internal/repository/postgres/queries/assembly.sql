-- name: CreateCar :exec
INSERT INTO cars (vin, brand, year) VALUES ($1, $2, $3);

-- name: CreateEngine :one
INSERT INTO engines (id, horsepower) VALUES ($1, $2) RETURNING id, horsepower;

-- name: CreateTransmission :one
INSERT INTO transmissions (id, type) VALUES ($1, $2) RETURNING id, type;

-- name: GetEngine :one
SELECT id, horsepower FROM engines WHERE id = $1;

-- name: GetTransmission :one
SELECT id, type FROM transmissions WHERE id = $1;

-- name: LinkComponents :exec
UPDATE cars SET engine_id = $2, transmission_id = $3 WHERE vin = $1;

-- name: GetCarSpec :one
SELECT
    c.vin, c.brand, c.year,
    e.id AS engine_id, e.horsepower,
    t.id AS trans_id, t.type AS trans_type
FROM cars c
LEFT JOIN engines e ON c.engine_id = e.id
LEFT JOIN transmissions t ON c.transmission_id = t.id
WHERE c.vin = $1;