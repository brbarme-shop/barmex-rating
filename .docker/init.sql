CREATE DATABASE rating_db;

\c rating_db;

-- Drop tables
DROP TABLE public.ratings;
DROP TABLE public.ratings_average;

CREATE TABLE public.ratings (
	rating_id serial4 NOT NULL,
	rating_item varchar NOT NULL,
	rating_avg numeric NOT NULL,
	CONSTRAINT ratings_pkey PRIMARY KEY (rating_id)
);

-- Permissions
ALTER TABLE public.ratings OWNER TO rating_user;
GRANT ALL ON TABLE public.ratings TO rating_user;

CREATE TABLE public.ratings_average (
	rating_id serial4 NOT NULL,
	rating_average_overall_rating int4 NOT NULL,
	rating_average_ratings int4 NOT NULL,
	CONSTRAINT ratings_average_pkey PRIMARY KEY (rating_id)
);

-- Permissions

ALTER TABLE public.ratings_average OWNER TO rating_user;
GRANT ALL ON TABLE public.ratings_average TO rating_user;


-- public.ratings_average foreign keys
ALTER TABLE public.ratings_average ADD CONSTRAINT fk_rating_id FOREIGN KEY (rating_id) REFERENCES public.ratings(rating_id);