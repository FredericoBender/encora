-- public.jobs definition

-- Drop table

-- DROP TABLE public.jobs;

CREATE TABLE public.jobs (
	job_id serial NOT NULL,
	title varchar NOT NULL,
	area varchar NOT NULL,
	country varchar NOT NULL,
	url text NOT NULL,
	description text NOT NULL DEFAULT ''::text,
	CONSTRAINT jobs_pk PRIMARY KEY (job_id)
);