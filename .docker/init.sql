CREATE DATABASE rating_db;

CREATE TABLE public.ratings (
	rating_id serial4 NOT NULL,
	rating_hash_id varchar (38) NOT null,
	rating_item_id varchar NOT NULL,
	rating_avg numeric NOT NULL,
	CONSTRAINT ratings_pkey PRIMARY KEY (rating_id),
	UNIQUE(rating_item_id, rating_hash_id)
);

CREATE TABLE public.ratings_stars (
	rating_star_id serial4 NOT NULL,
	rating_star int NOT NULL,
	CONSTRAINT ratings_stars_pkey PRIMARY KEY (rating_star_id)
);

CREATE TABLE public.ratings_averages (
	rating_id serial4 NOT NULL,
	rating_star_id serial4 NOT NULL,
	rating_count int NOT NULL
);
 
ALTER TABLE public.ratings_averages ADD CONSTRAINT fk_rating_id FOREIGN KEY (rating_id) REFERENCES public.ratings(rating_id);
ALTER TABLE public.ratings_stars ADD CONSTRAINT fk_rating_star_id FOREIGN KEY (rating_star_id) REFERENCES public.ratings_stars(rating_star_id);

INSERT INTO ratings (rating_hash_id, rating_item_id, rating_avg) 
VALUES('8925314f-3dc0-48b3-8a2e-2778350f28cf','0798112345321', 4.11);

INSERT INTO ratings_stars (rating_star) 
VALUES(5);