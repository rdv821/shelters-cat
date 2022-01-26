CREATE TABLE CATS
(
    id         uuid         NOT NULL UNIQUE,
    name       varchar(255) NOT NULL,
    age        int          NOT NULL,
    vaccinated boolean      NOT NULL default false
);
