CREATE TABLE IF NOT EXISTS phone (
    id uuid PRIMARY KEY,
    number varchar(30) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS phone_number_idx ON phone (number);
