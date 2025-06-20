CREATE TABLE hubs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL UNIQUE,
    address VARCHAR(512) NOT NULL,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(10),
    contact_name VARCHAR(255),
    contact_email VARCHAR(255)
);

CREATE TABLE skus (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    code VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    tenant_id VARCHAR(255) NOT NULL,
    seller_id VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    price DOUBLE PRECISION
);

CREATE INDEX idx_tenant_id ON skus(tenant_id);
CREATE INDEX idx_seller_id ON skus(seller_id);

CREATE TABLE inventories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    hub_id INT NOT NULL,
    sku_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    UNIQUE (hub_id, sku_id),
    FOREIGN KEY (hub_id) REFERENCES hubs(id) ON DELETE CASCADE,
    FOREIGN KEY (sku_id) REFERENCES skus(id) ON DELETE CASCADE
);
