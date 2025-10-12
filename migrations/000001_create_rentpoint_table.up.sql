CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.rentpoint (
   id uuid DEFAULT uuid_generate_v4() NOT NULL,
   name varchar(50) NOT NULL,
   addr varchar(50) NOT NULL
);