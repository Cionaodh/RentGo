CREATE TABLE public.templates (
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   name varchar(50) NOT NULL,
   description varchar(50) NOT NULL,
   price numeric(5, 2)
);