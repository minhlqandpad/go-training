CREATE TABLE tenants (
    id VARCHAR(50) PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tenants (id, Name, email) VALUES
('tenant-123', 'test', 'test@example.com'),
('tenant-456', 'alice', 'alice@example.com');
