BEGIN;

-- Create applications table
CREATE TABLE applications (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    link VARCHAR(500) NOT NULL,
    show_in_jumbotron BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER
);

COMMIT;
