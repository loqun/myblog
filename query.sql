-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetAllBlog :many
SELECT b.id, b.title, b.content, b.tag_id, b.author, b.created_at, t.name as tag_name
FROM blogs b
LEFT JOIN tags t ON b.tag_id = t.id
ORDER BY b.created_at DESC;

-- name: GetBlogsPaginated :many
SELECT b.id, b.title, b.content, b.tag_id, b.author, b.created_at, t.name as tag_name
FROM blogs b
LEFT JOIN tags t ON b.tag_id = t.id
ORDER BY b.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountBlogs :one
SELECT COUNT(*) FROM blogs;

-- name: GetBlogByID :one
SELECT b.id, b.title, b.content, b.tag_id, b.author, b.created_at, t.name as tag_name
FROM blogs b
LEFT JOIN tags t ON b.tag_id = t.id
WHERE b.id = ?;

-- name: GetOrCreateTag :one
INSERT INTO tags (id, name) VALUES (?, ?) 
ON CONFLICT(name) DO UPDATE SET name = name
RETURNING id, name;

-- name: StoreBlog :one
INSERT INTO blogs (id, title, content, tag_id, author, created_at) VALUES (?, ?, ?, ?, ?, ?) RETURNING id, title, content, tag_id, author, created_at;