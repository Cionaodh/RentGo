CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.rentpoints (
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   "name" varchar(50) NOT NULL,
   addr varchar(50) NOT NULL,
   CONSTRAINT rentpoints_pkey PRIMARY KEY (id),
   CONSTRAINT rentpoints_name_key UNIQUE ("name"),
   CONSTRAINT rentpoints_addr_key UNIQUE (addr)
   -- TODO: добавить столбец is_deleted - чтобы не удалять записи, а просто помечать, что она удалена 
);

CREATE TABLE public.templates (
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   "name" varchar(50) NOT NULL,
   "description" varchar(50) NOT NULL,
   price numeric(7, 2) NOT NULL,
   CONSTRAINT templates_pkey PRIMARY KEY (id),
   CONSTRAINT templates_name_key UNIQUE ("name")
   -- TODO: добавить столбец is_deleted
);

CREATE TYPE status_product_type AS ENUM (
    'Unused', -- неиспользуется
    'Free', -- свободен (прикреплен к точке проката и готов к аренде)
    'Reserved', -- временно зарезервирован 
    'Rented', -- арендован
    'Retired' -- снят с эксплуатации
);
   
CREATE TABLE public.products(
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   "status" status_product_type NOT NULL,
   template_id uuid NOT NULL, -- TODO: шаблон у продукта менять нельзя
   rentpoint_id uuid NULL,

   CONSTRAINT products_pkey PRIMARY KEY (id),
   CONSTRAINT product_template_fkey FOREIGN KEY (template_id) 
      REFERENCES public.templates(id) 
      ON DELETE NO ACTION
      ON UPDATE CASCADE,
   CONSTRAINT product_rentpoint_fkey FOREIGN KEY (rentpoint_id) 
      REFERENCES public.rentpoints(id) 
      ON DELETE SET NULL 
      ON UPDATE CASCADE
);

CREATE TYPE status_order_type AS ENUM (
   'Active', -- Активный
   'Completed',-- Завершенный
   'Expired'-- Просрочен - если пользователь не завершил поездку за купленное время
);

CREATE TABLE public.orders(
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   "status" status_order_type NOT NULL,
   product_id uuid NOT NULL,-- id продукта (внешний ключ)
   starting_point_id uuid NOT NULL,-- id начальной точки проката (внешний ключ)
   finishing_point_id uuid NULL,-- id конечной точки проката (внешний ключ)
   rent_started_at timestamptz NOT NULL,-- время начала аренды 
   rent_finished_at timestamptz NULL,-- время завершения аренды

   CONSTRAINT orders_pkey PRIMARY KEY (id),
   CONSTRAINT order_product_fkey FOREIGN KEY (product_id)
      REFERENCES public.products(id)
      ON DELETE NO ACTION -- TODO: по сути если мы удалим продукт из базы данных - мы потеряем информацию о заказе - о продукте который был арендован, так как мы не найдем шаблон продукта
      ON UPDATE CASCADE,
   CONSTRAINT order_starting_point_fkey FOREIGN KEY (starting_point_id)
      REFERENCES public.rentpoints(id)
      ON DELETE NO ACTION
      ON UPDATE CASCADE,
   CONSTRAINT order_finishing_point_fkey FOREIGN KEY (finishing_point_id)
      REFERENCES public.rentpoints(id)
      ON DELETE NO ACTION
      ON UPDATE CASCADE
);

--------------------------------------
-- ИНДЕКСЫ
--------------------------------------

CREATE INDEX idx_products_template_id ON public.products(template_id);
CREATE INDEX idx_products_rentpoint_id ON public.products(rentpoint_id);
CREATE INDEX idx_orders_product_id ON public.orders(product_id);
CREATE INDEX idx_orders_starting_point_id ON public.orders(starting_point_id);
CREATE INDEX idx_orders_finishing_point_id ON public.orders(finishing_point_id);
CREATE INDEX idx_orders_status ON public.orders(status);
CREATE INDEX idx_orders_dates ON public.orders(rent_started_at, rent_finished_at);

--------------------------------------
-- ТЕСТОВОЕ ЗАПОЛНЕНИЕ
--------------------------------------

INSERT INTO public.rentpoints (id, "name", addr) VALUES
    ('550e8400-e29b-41d4-a716-446655440001'::uuid, 'Центральный вокзал', 'ул. Ленина, 1'),
    ('550e8400-e29b-41d4-a716-446655440002'::uuid, 'Торговый центр "Молл"', 'пр. Мира, 45'),
    ('550e8400-e29b-41d4-a716-446655440003'::uuid, 'Парк культуры', 'ул. Парковая, 12'),
    ('550e8400-e29b-41d4-a716-446655440004'::uuid, 'Университет', 'ул. Студенческая, 25'),
    ('550e8400-e29b-41d4-a716-446655440005'::uuid, 'Жилой район "Северный"', 'ул. Северная, 8');

INSERT INTO public.templates (id, "name", description, price) VALUES
    ('550e8400-e29b-41d4-a716-446655550001'::uuid, 'Эконом', 'Базовая модель для города', 100.00),
    ('550e8400-e29b-41d4-a716-446655550002'::uuid, 'Стандарт', 'Средний класс с корзиной', 150.00),
    ('550e8400-e29b-41d4-a716-446655550003'::uuid, 'Премиум', 'Электрический, с GPS', 250.00),
    ('550e8400-e29b-41d4-a716-446655550004'::uuid, 'Грузовой', 'Для перевозки небольших грузов', 200.00),
    ('550e8400-e29b-41d4-a716-446655550005'::uuid, 'Спорт', 'Для активной езды', 180.00);

INSERT INTO public.products (id, "status", template_id, rentpoint_id) VALUES
    ('550e8400-e29b-41d4-a716-446655660001'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550001'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid),
    ('550e8400-e29b-41d4-a716-446655660002'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550001'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid),
    ('550e8400-e29b-41d4-a716-446655660003'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550002'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid),
    ('550e8400-e29b-41d4-a716-446655660004'::uuid, 'Rented', '550e8400-e29b-41d4-a716-446655550002'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid),
    ('550e8400-e29b-41d4-a716-446655660005'::uuid, 'Reserved', '550e8400-e29b-41d4-a716-446655550003'::uuid, '550e8400-e29b-41d4-a716-446655440003'::uuid),
    ('550e8400-e29b-41d4-a716-446655660006'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550003'::uuid, '550e8400-e29b-41d4-a716-446655440003'::uuid),
    ('550e8400-e29b-41d4-a716-446655660007'::uuid, 'Unused', '550e8400-e29b-41d4-a716-446655550004'::uuid, NULL),
    ('550e8400-e29b-41d4-a716-446655660008'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550005'::uuid, '550e8400-e29b-41d4-a716-446655440004'::uuid),
    ('550e8400-e29b-41d4-a716-446655660009'::uuid, 'Free', '550e8400-e29b-41d4-a716-446655550005'::uuid, '550e8400-e29b-41d4-a716-446655440005'::uuid),
    ('550e8400-e29b-41d4-a716-446655660010'::uuid, 'Retired', '550e8400-e29b-41d4-a716-446655550001'::uuid, NULL);

INSERT INTO public.orders (id, "status", product_id, starting_point_id, finishing_point_id, rent_started_at, rent_finished_at) VALUES
    -- Активный заказ (самокат в аренде)
    ('550e8400-e29b-41d4-a716-446655770001'::uuid, 'Active', 
     '550e8400-e29b-41d4-a716-446655660004'::uuid, 
     '550e8400-e29b-41d4-a716-446655440002'::uuid, 
     NULL,
     '2024-01-15 14:30:00+03',
     NULL),
    
    -- Завершенный заказ (успешно завершен)
    ('550e8400-e29b-41d4-a716-446655770002'::uuid, 'Completed', 
     '550e8400-e29b-41d4-a716-446655660003'::uuid, 
     '550e8400-e29b-41d4-a716-446655440002'::uuid, 
     '550e8400-e29b-41d4-a716-446655440004'::uuid,
     '2024-01-15 10:00:00+03',
     '2024-01-15 11:30:00+03'),
    
    -- Просроченный заказ (не завершен вовремя)
    ('550e8400-e29b-41d4-a716-446655770003'::uuid, 'Expired', 
     '550e8400-e29b-41d4-a716-446655660005'::uuid, 
     '550e8400-e29b-41d4-a716-446655440003'::uuid, 
     NULL,
     '2024-01-14 18:00:00+03',
     NULL),
    
    -- Завершенный заказ (аренда на день)
    ('550e8400-e29b-41d4-a716-446655770004'::uuid, 'Completed', 
     '550e8400-e29b-41d4-a716-446655660001'::uuid, 
     '550e8400-e29b-41d4-a716-446655440001'::uuid, 
     '550e8400-e29b-41d4-a716-446655440005'::uuid,
     '2024-01-13 09:00:00+03',
     '2024-01-13 19:00:00+03'),
    
    -- Активный заказ (только начат)
    ('550e8400-e29b-41d4-a716-446655770005'::uuid, 'Active', 
     '550e8400-e29b-41d4-a716-446655660008'::uuid, 
     '550e8400-e29b-41d4-a716-446655440004'::uuid, 
     NULL,
     '2024-01-15 15:45:00+03',
     NULL),
    
    -- Завершенный заказ (короткая поездка)
    ('550e8400-e29b-41d4-a716-446655770006'::uuid, 'Completed', 
     '550e8400-e29b-41d4-a716-446655660006'::uuid, 
     '550e8400-e29b-41d4-a716-446655440003'::uuid, 
     '550e8400-e29b-41d4-a716-446655440001'::uuid,
     '2024-01-14 12:15:00+03',
     '2024-01-14 13:45:00+03'),
    
    -- Просроченный заказ (давно просрочен)
    ('550e8400-e29b-41d4-a716-446655770007'::uuid, 'Expired', 
     '550e8400-e29b-41d4-a716-446655660002'::uuid, 
     '550e8400-e29b-41d4-a716-446655440001'::uuid, 
     NULL,
     '2024-01-10 08:00:00+03',
     NULL),
    
    -- Завершенный заказ (возврат в ту же точку)
    ('550e8400-e29b-41d4-a716-446655770008'::uuid, 'Completed', 
     '550e8400-e29b-41d4-a716-446655660009'::uuid, 
     '550e8400-e29b-41d4-a716-446655440005'::uuid, 
     '550e8400-e29b-41d4-a716-446655440005'::uuid,
     '2024-01-12 16:30:00+03',
     '2024-01-12 18:00:00+03');