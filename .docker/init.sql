CREATE DATABASE rating_db;

CREATE TABLE ratings_avarages (
	rating_id serial4 NOT NULL,
	rating_hash_id VARCHAR (38) NOT NULL,
	rating_item_id VARCHAR NOT null,
	rating_avg NUMERIC NOT NULL,
	rating_start_i INTEGER DEFAULT 1 NOT NULL,
	rating_start_i_count NUMERIC DEFAULT 0 NOT NULL,
	rating_start_ii INTEGER DEFAULT 2 NOT NULL,
	rating_start_ii_count NUMERIC DEFAULT 0 NOT NULL,
	rating_start_iii INTEGER DEFAULT 3 NOT NULL,
	rating_start_iii_count NUMERIC DEFAULT 0 NOT NULL,
	rating_start_iv INTEGER DEFAULT 4 NOT NULL,
	rating_start_iv_count NUMERIC DEFAULT 0 NOT NULL,
	rating_start_x INTEGER DEFAULT 5 NOT NULL,
	rating_start_x_count NUMERIC DEFAULT 0 NOT NULL,
	CONSTRAINT ratings_pkey PRIMARY KEY (rating_id),
	UNIQUE(rating_item_id, rating_hash_id)
);