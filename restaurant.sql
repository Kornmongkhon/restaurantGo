-- ลบฐานข้อมูลเดิมถ้ามีอยู่
DROP DATABASE IF EXISTS restaurant;

-- สร้างฐานข้อมูลใหม่
CREATE DATABASE restaurant;
USE restaurant;

-- ลบตาราง tables (โต๊ะ) ถ้ามีอยู่
DROP TABLE IF EXISTS tables;

-- สร้างตาราง tables (โต๊ะ)
CREATE TABLE tables (
                        id INT AUTO_INCREMENT PRIMARY KEY,
                        table_number INT UNIQUE NOT NULL,
                        is_occupied BOOLEAN DEFAULT FALSE,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ลบตาราง menu_items (เมนูอาหาร) ถ้ามีอยู่
DROP TABLE IF EXISTS menu_items;

-- สร้างตาราง menu_items (เมนูอาหาร)
CREATE TABLE menu_items (
                            id INT AUTO_INCREMENT PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            description TEXT,
                            price DECIMAL(10, 2) NOT NULL,
                            is_available BOOLEAN DEFAULT TRUE,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ลบตาราง orders (ออเดอร์) ถ้ามีอยู่
DROP TABLE IF EXISTS orders;

-- สร้างตาราง orders (ออเดอร์)
CREATE TABLE orders (
                        id INT AUTO_INCREMENT PRIMARY KEY,
                        table_id INT,
                        total_amount DECIMAL(10, 2) DEFAULT 0,
                        status ENUM('created', 'updated', 'canceled', 'completed') DEFAULT 'created',
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE CASCADE
);

-- ลบตาราง order_items (รายการออเดอร์) ถ้ามีอยู่
DROP TABLE IF EXISTS order_items;

-- สร้างตาราง order_items (รายการออเดอร์)
CREATE TABLE order_items (
                             id INT AUTO_INCREMENT PRIMARY KEY,
                             order_id INT,
                             menu_item_id INT,
                             quantity INT NOT NULL,
                             price DECIMAL(10, 2) NOT NULL,
                             FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
                             FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE CASCADE
);

-- ลบตาราง bills (บิล) ถ้ามีอยู่
DROP TABLE IF EXISTS bills;

-- สร้างตาราง bills (บิล)
CREATE TABLE bills (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       order_id INT,
                       table_id INT,
                       total_amount DECIMAL(10, 2) NOT NULL,
                       bill_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
                       FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE CASCADE
);

-- ลบตาราง reviews (รีวิว) ถ้ามีอยู่
DROP TABLE IF EXISTS reviews;

-- สร้างตาราง reviews (รีวิว)
CREATE TABLE reviews (
                         id INT AUTO_INCREMENT PRIMARY KEY,
                         menu_item_id INT,
                         order_id INT,
                         rating INT CHECK (rating >= 1 AND rating <= 5),
                         comment TEXT,
                         review_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE CASCADE,
                         FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- ข้อมูลตัวอย่างสำหรับตาราง tables
INSERT INTO tables (table_number, is_occupied) VALUES
                                                   (1, FALSE),
                                                   (2, FALSE),
                                                   (3, FALSE),
                                                   (4, FALSE),
                                                   (5, FALSE);

-- ข้อมูลตัวอย่างสำหรับตาราง menu_items (ราคาจะเป็นบาท)
INSERT INTO menu_items (name, description, price, is_available) VALUES
                                                                    ('Spaghetti Carbonara', 'Classic Italian pasta dish with creamy sauce. ราคา: 299 บาท', 299.00, TRUE),
                                                                    ('Margherita Pizza', 'Simple yet delicious pizza with fresh mozzarella. ราคา: 259 บาท', 259.00, TRUE),
                                                                    ('Caesar Salad', 'Crisp romaine lettuce with Caesar dressing and croutons. ราคา: 179 บาท', 179.00, TRUE),
                                                                    ('Tiramisu', 'Delicious coffee-flavored Italian dessert. ราคา: 129 บาท', 129.00, TRUE);
