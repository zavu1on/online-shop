CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   login varchar(100) UNIQUE NOT NULL,
   password varchar(300) NOT NULL
);

CREATE TABLE IF NOT EXISTS baskets(
   id serial PRIMARY KEY,
   user_id int REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS goods(
   id serial PRIMARY KEY,
   title varchar(100) UNIQUE NOT NULL,
   description text NOT NULL,
   price int NOT NULL,
   count int NOT NULL 
);

CREATE TABLE IF NOT EXISTS baskets_goods(
  good_id int REFERENCES goods(id),
  basket_id int REFERENCES baskets(id),

  CONSTRAINT baskets_goods_pkey PRIMARY KEY (good_id, basket_id)
)