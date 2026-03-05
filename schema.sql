
CREATE TABLE IF NOT EXISTS users (   
    id TEXT PRIMARY KEY,   
    name TEXT NOT NULL,   
    email TEXT NOT NULL UNIQUE  
);      


-- // tags table
CREATE TABLE IF NOT EXISTS tags (   
    id TEXT PRIMARY KEY,   
    name TEXT NOT NULL UNIQUE  
    
);

CREATE TABLE IF NOT EXISTS blogs (   
    id TEXT PRIMARY KEY,   
    title TEXT NOT NULL,   
    content TEXT NOT NULL,
    tag_id TEXT REFERENCES tags(id),
    author TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);   






