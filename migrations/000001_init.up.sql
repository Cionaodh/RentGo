CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.rentpoints (
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   "name" varchar(50) NOT NULL,
   addr varchar(50) NOT NULL,
   CONSTRAINT rentpoint_pkey PRIMARY KEY (id),
   CONSTRAINT rentpoint_name_key UNIQUE ("name"),
   CONSTRAINT rentpoint_addr_key UNIQUE (addr)
);

CREATE TABLE public.templates (
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   "name" varchar(50) NOT NULL,
   description varchar(50) NOT NULL,
   price numeric(10, 2) NOT NULL,
   CONSTRAINT templates_pkey PRIMARY KEY (id)
);

CREATE TYPE status_type AS ENUM (
    'Unused', 
    'Free',
    'Reserved',
    'Rented'
);

CREATE TABLE public.products(
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   "status" status_type NOT NULL,
   template_id uuid NOT NULL,
   rentpoint_id uuid NULL,
   CONSTRAINT products_pkey PRIMARY KEY (id),
   CONSTRAINT product_template_fkey FOREIGN KEY (template_id) 
      REFERENCES public.templates(id) 
      ON DELETE CASCADE 
      ON UPDATE CASCADE,
   CONSTRAINT product_rentpoint_fkey FOREIGN KEY (rentpoint_id) 
      REFERENCES public.rentpoints(id) 
      ON DELETE SET NULL 
      ON UPDATE CASCADE -- Добавлено действие для ON UPDATE
);

-- Заполнение начальными значениями
INSERT INTO public.rentpoints (id, "name", addr) VALUES
   ('550e8400-e29b-41d4-a716-446655440001'::uuid, 'rentpoint1', 'First1'),
   ('550e8400-e29b-41d4-a716-446655440002'::uuid, 'rentpoint2', 'First2'),
   ('550e8400-e29b-41d4-a716-446655440003'::uuid, 'rentpoint3', 'First3');

INSERT INTO public.templates (id, "name", description, price) VALUES 
   ('550e8400-e29b-41d4-a716-446655550001'::uuid, 'template1', 'description1', 1500), -- Исправлены кавычки и опечатка
   ('550e8400-e29b-41d4-a716-446655550002'::uuid, 'template2', 'description2', 2000),
   ('550e8400-e29b-41d4-a716-446655550003'::uuid, 'template3', 'description3', 1000),
   ('550e8400-e29b-41d4-a716-446655550004'::uuid, 'template4', 'description4', 3000);

-- Пример корректной вставки в products (при необходимости):
-- INSERT INTO public.products (id, status, template_id, rentpoint_id) VALUES
--    (uuid_generate_v4(), 'Free', '550e8400-e29b-41d4-a716-446655550001'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid);