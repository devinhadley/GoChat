CREATE TABLE users (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username varchar(100) NOT NULL,
    password_hash varchar(255) NOT NULL,
    sign_up_date date NOT NULL DEFAULT CURRENT_DATE,
    is_active boolean NOT NULL DEFAULT true
);
