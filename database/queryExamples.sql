select title,area,country,count(title) as total from jobs group by title,area,country order by total desc; --only show how many instance we have grouped by title,area and country

select count(1) from jobs; --show the number of jobs in tha table named 'jobs'

select distinct title from jobs order by title; --show only the different jobs from the database

select distinct country from jobs order by country ; --only shows the different countries that have encora jobs

select distinct area from jobs order by area ;  --only shows the different areas of software that have work on the encora website