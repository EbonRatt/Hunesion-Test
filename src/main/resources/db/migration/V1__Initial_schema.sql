-- Create users table with UUID
CREATE TABLE IF NOT EXISTS users (
                                     id BINARY(16) PRIMARY KEY,
                                     username VARCHAR(50) NOT NULL UNIQUE,
                                     email VARCHAR(100) NOT NULL UNIQUE,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Create inventories table with UUID
CREATE TABLE IF NOT EXISTS inventories (
                                           id BINARY(16) PRIMARY KEY,
                                           user_id BINARY(16) NOT NULL UNIQUE,
                                           capacity INT NOT NULL DEFAULT 20,
                                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Create items table with UUID
CREATE TABLE IF NOT EXISTS items (
                                     id BINARY(16) PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
                                     description TEXT,
                                     value DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
                                     rarity ENUM('COMMON', 'UNCOMMON', 'RARE', 'EPIC', 'LEGENDARY') NOT NULL DEFAULT 'COMMON',
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Create inventory_items join table with UUID foreign keys
CREATE TABLE IF NOT EXISTS inventory_items (
                                               inventory_id BINARY(16) NOT NULL,
                                               item_id BINARY(16) NOT NULL,
                                               quantity INT NOT NULL DEFAULT 1,
                                               acquired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               PRIMARY KEY (inventory_id, item_id),
                                               FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE CASCADE,
                                               FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Create indexes for better query performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_inventories_user_id ON inventories(user_id);
CREATE INDEX idx_items_name ON items(name);
CREATE INDEX idx_inventory_items_inventory ON inventory_items(inventory_id);
CREATE INDEX idx_inventory_items_item ON inventory_items(item_id);