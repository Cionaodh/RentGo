CREATE TABLE users (
   id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
   email varchar(255) UNIQUE NOT NULL,
   username varchar(255) NOT NULL,
   password_hash varchar(255) NOT NULL,
   created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE rentpoints (
   id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
   "name" varchar(50) UNIQUE NOT NULL,
   addr varchar(50) UNIQUE NOT NULL,
   is_deleted boolean NOT NULL DEFAULT false
);

CREATE TABLE templates (
   id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
   "name" varchar(50) UNIQUE NOT NULL,
   "description" varchar(50) NOT NULL,
   price numeric(7, 2) NOT NULL,
   is_deleted boolean NOT NULL DEFAULT false
);

CREATE TYPE status_product_type AS ENUM (
    'Unused', -- неиспользуется
    'Free', -- свободен (прикреплен к точке проката и готов к аренде)
    'Reserved', -- временно зарезервирован 
    'Rented', -- арендован
    'Retired' -- снят с эксплуатации
);
   
CREATE TABLE products(
   id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
   "status" status_product_type NOT NULL,
   template_id uuid NOT NULL REFERENCES templates(id) ON DELETE NO ACTION ON UPDATE CASCADE,
   rentpoint_id uuid REFERENCES rentpoints(id) ON DELETE SET NULL ON UPDATE CASCADE
);

-- Статусы заказа
CREATE TYPE status_order_type AS ENUM (
   'Active',     -- Активный заказ
   'Completed',  -- Завершён
   'Expired'     -- Просрочен
);

CREATE TABLE orders (
   id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
   user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   "status" status_order_type NOT NULL DEFAULT 'Active',
   
   product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT ON UPDATE CASCADE,
   starting_point_id uuid NOT NULL REFERENCES rentpoints(id) ON DELETE NO ACTION ON UPDATE CASCADE,
   finishing_point_id uuid REFERENCES rentpoints(id) ON DELETE NO ACTION ON UPDATE CASCADE,

   rent_started_at timestamptz NOT NULL DEFAULT now(),
   rent_finished_at timestamptz,

   CONSTRAINT order_rent_time_check CHECK (
      rent_finished_at IS NULL OR rent_finished_at >= rent_started_at
   )
);

--------------------------------------
-- ИНДЕКСЫ 
--------------------------------------
CREATE INDEX idx_products_template_id ON products(template_id);
CREATE INDEX idx_products_rentpoint_id ON products(rentpoint_id);

CREATE INDEX idx_orders_user_id ON orders(user_id); 
CREATE INDEX idx_orders_product_id ON orders(product_id);
CREATE INDEX idx_orders_starting_point_id ON orders(starting_point_id);
CREATE INDEX idx_orders_finishing_point_id ON orders(finishing_point_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_dates ON orders(rent_started_at, rent_finished_at);

--------------------------------------
-- ТЕСТОВОЕ ЗАПОЛНЕНИЕ
--------------------------------------


INSERT INTO users (id, email, username, password_hash) VALUES
    ('550e8400-e29b-41d4-a716-111122220001'::uuid, 'ivanov', 'Иван Иванов', '$2a$10$somehash1...'),
    ('550e8400-e29b-41d4-a716-111122220002'::uuid, 'petrov', 'Петр Петров', '$2a$10$somehash2...');


INSERT INTO rentpoints (id, "name", addr) VALUES
    ('550e8400-e29b-41d4-a716-446655440001'::uuid, 'Точка А', 'ул. Ленина, 1'),
    ('550e8400-e29b-41d4-a716-446655440002'::uuid, 'Точка B', 'пр. Мира, 45'),
    ('550e8400-e29b-41d4-a716-446655440003'::uuid, 'Точка C', 'ул. Парковая, 12'),
    ('550e8400-e29b-41d4-a716-446655440004'::uuid, 'Точка D', 'ул. Студенческая, 25'),
    ('550e8400-e29b-41d4-a716-446655440005'::uuid, 'Точка E', 'ул. Северная, 8');


INSERT INTO templates (id, "name", description, price) VALUES
    ('550e8400-e29b-41d4-a716-446655550001'::uuid, 'Велосипед', 'Базовая модель для города', 100.00),
    ('550e8400-e29b-41d4-a716-446655550002'::uuid, 'Самокат', 'Средний класс с корзиной', 150.00),
    ('550e8400-e29b-41d4-a716-446655550003'::uuid, 'Ролики', 'Электрический, с GPS', 250.00),
    ('550e8400-e29b-41d4-a716-446655550004'::uuid, 'Самолет', 'Для перевозки небольших грузов', 200.00),
    ('550e8400-e29b-41d4-a716-446655550005'::uuid, 'Хеликоптер', 'Для активной езды', 180.00);


INSERT INTO products (id, "status", template_id, rentpoint_id) VALUES
    ('550e8400-e29b-41d4-a716-446655660001'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550001'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid),
    ('550e8400-e29b-41d4-a716-446655660002'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550001'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid),
    ('550e8400-e29b-41d4-a716-446655660003'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550002'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid),
    ('550e8400-e29b-41d4-a716-446655660004'::uuid, 'Rented', '550e8400-e29b-41d4-a716-446655550002'::uuid, NULL),
    ('550e8400-e29b-41d4-a716-446655660005'::uuid, 'Reserved', '550e8400-e29b-41d4-a716-446655550003'::uuid, '550e8400-e29b-41d4-a716-446655440003'::uuid),
    ('550e8400-e29b-41d4-a716-446655660006'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550003'::uuid, '550e8400-e29b-41d4-a716-446655440003'::uuid),
    ('550e8400-e29b-41d4-a716-446655660007'::uuid, 'Unused', '550e8400-e29b-41d4-a716-446655550004'::uuid, NULL),
    ('550e8400-e29b-41d4-a716-446655660008'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550005'::uuid, '550e8400-e29b-41d4-a716-446655440004'::uuid),
    ('550e8400-e29b-41d4-a716-446655660009'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550005'::uuid, '550e8400-e29b-41d4-a716-446655440005'::uuid),
    ('550e8400-e29b-41d4-a716-446655660010'::uuid, 'Retired', '550e8400-e29b-41d4-a716-446655550001'::uuid, NULL);


INSERT INTO orders (id, user_id, "status", product_id, starting_point_id, finishing_point_id, rent_started_at, rent_finished_at) VALUES
    -- Завершённый заказ
    ('550e8400-e29b-41d4-a716-777788880001'::uuid, '550e8400-e29b-41d4-a716-111122220001'::uuid, 'Completed', '550e8400-e29b-41d4-a716-446655660003'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid, now() - interval '2 days', now() - interval '1 day'),

    -- Текущий активный заказ
    ('550e8400-e29b-41d4-a716-777788880002'::uuid, '550e8400-e29b-41d4-a716-111122220002'::uuid, 'Active', '550e8400-e29b-41d4-a716-446655660004'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid, NULL, now() - interval '2 hours', NULL);