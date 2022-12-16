
DROP TABLE if exists person;

-- create the person table

CREATE TABLE person (
	first_name varchar(50),
	last_name varchar(50),
	dob date
);

-- load table

\copy person FROM 'person.csv' DELIMITER ',' CSV HEADER;
