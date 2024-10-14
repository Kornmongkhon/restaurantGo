-- ลบฐานข้อมูลเดิมถ้ามีอยู่
DROP DATABASE IF EXISTS restaurant;

-- สร้างฐานข้อมูลใหม่
CREATE DATABASE restaurant;
USE restaurant;

-- ลบตาราง tables (โต๊ะ) ถ้ามีอยู่
DROP TABLE IF EXISTS tables;

-- สร้างตาราง tables (โต๊ะ)
CREATE TABLE tables (
                        table_id INT AUTO_INCREMENT PRIMARY KEY,
                        table_number INT UNIQUE NOT NULL,
                        table_status ENUM('available', 'occupied') DEFAULT 'available',
                        is_deleted BOOLEAN DEFAULT FALSE,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NULL DEFAULT NULL
);

-- ลบตาราง menu_items (เมนูอาหาร) ถ้ามีอยู่
DROP TABLE IF EXISTS menu_items;

-- สร้างตาราง menu_items (เมนูอาหาร)
CREATE TABLE menu_items (
                            menu_items_id INT AUTO_INCREMENT PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            description TEXT,
                            price DECIMAL(10, 2) NOT NULL,
                            file_path VARCHAR(255),
                            is_available BOOLEAN DEFAULT TRUE,
                            is_deleted BOOLEAN DEFAULT FALSE,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ลบตาราง orders (ออเดอร์) ถ้ามีอยู่
DROP TABLE IF EXISTS orders;

-- สร้างตาราง orders (ออเดอร์)
CREATE TABLE orders (
                        order_id INT AUTO_INCREMENT PRIMARY KEY,
                        table_id INT,
                        status ENUM('created', 'prepare', 'canceled', 'completed', 'paid') DEFAULT 'created',
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NULL DEFAULT NULL,
                        is_deleted BOOLEAN DEFAULT FALSE,
                        FOREIGN KEY (table_id) REFERENCES tables(table_id) ON DELETE CASCADE
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
                             FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
                             FOREIGN KEY (menu_item_id) REFERENCES menu_items(menu_items_id) ON DELETE CASCADE
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
                       FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
                       FOREIGN KEY (table_id) REFERENCES tables(table_id) ON DELETE CASCADE
);

-- ลบตาราง reviews (รีวิว) ถ้ามีอยู่
DROP TABLE IF EXISTS reviews;

-- สร้างตาราง reviews (รีวิว)
CREATE TABLE reviews (
                         id INT AUTO_INCREMENT PRIMARY KEY,
#                          menu_item_id INT,
                         order_id INT,
                         rating INT CHECK (rating >= 1 AND rating <= 5),
                         comment TEXT,
                         is_deleted BOOLEAN DEFAULT FALSE,
                         review_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
#                          FOREIGN KEY (menu_item_id) REFERENCES menu_items(menu_items_id) ON DELETE CASCADE,
                         FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);

-- ข้อมูลตัวอย่างสำหรับตาราง tables
INSERT INTO tables (table_number, is_deleted) VALUES
                                                   (1, FALSE),
                                                   (2, FALSE),
                                                   (3, FALSE),
                                                   (4, FALSE),
                                                   (5, FALSE);

-- ข้อมูลตัวอย่างสำหรับตาราง menu_items (ราคาจะเป็นบาท)
-- เพิ่มเมนูใหม่ 15 รายการ
INSERT INTO menu_items (name, description, price, file_path, is_available, is_deleted) VALUES
                                                                    ('Spaghetti Carbonara', 'Classic Italian pasta with creamy sauce', 150.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Spaghetti Carbonara.jpg',true, false),
                                                                    ('Margherita Pizza', 'Traditional pizza with tomato, mozzarella, and basil', 200.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Margherita Pizza.jpg', true, false),
                                                                    ('Caesar Salad', 'Crispy romaine lettuce with Caesar dressing', 120.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Caesar Salad.jpg',true, false),
                                                                    ('Grilled Salmon', 'Fresh salmon grilled with herbs and lemon', 350.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Grilled Salmon.jpg',true, false),
                                                                    ('Chicken Parmesan', 'Crispy chicken breast with marinara and mozzarella', 250.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Chicken Parmesan.jpg',true, false),
                                                                    ('Beef Burger', 'Juicy beef patty with cheese and lettuce', 180.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Beef Burger.jpg',true, false),
                                                                    ('French Fries', 'Golden and crispy fries', 80.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\French Fries.jpg',true,false),
                                                                    ('Vegetable Stir Fry', 'Mixed vegetables stir-fried with soy sauce', 140.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Vegetable Stir Fry.jpg',true, false),
                                                                    ('Pad Thai', 'Classic Thai stir-fried noodles with shrimp', 150.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Pad Thai.jpg',true, false),
                                                                    ('Tom Yum Soup', 'Spicy and sour Thai soup with shrimp', 180.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Tom Yum Soup.jpg',true, false),
                                                                    ('Chicken Tikka Masala', 'Spicy chicken in creamy tomato sauce', 220.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Chicken Tikka Masala.jpg',true, false),
                                                                    ('Sushi Platter', 'Assorted sushi with fresh fish and vegetables', 300.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Sushi Platter.jpg',true, false),
                                                                    ('Ramen', 'Japanese noodle soup with pork and egg', 180.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Ramen.jpg',true, false),
                                                                    ('Pancakes', 'Fluffy pancakes with syrup and butter', 100.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Pancakes.jpg',true, false),
                                                                    ('Chocolate Cake', 'Rich chocolate cake with fudge icing', 90.00, 'K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\Chocolate Cake.jpg',true, false);

-- ทำการอัปเดตค่า is_available ให้เป็น false สำหรับ 5 เมนูแบบสุ่ม
UPDATE menu_items SET is_available = false WHERE name IN (
                                                          'Caesar Salad',
                                                          'French Fries',
                                                          'Sushi Platter',
                                                          'Tom Yum Soup',
                                                          'Pancakes'
    );
