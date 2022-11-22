CREATE TABLE IF NOT EXISTS `users`
(
    `id`       BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(50) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL ,
    `role`     VARCHAR(15) NOT NULL
);

INSERT INTO users (username, password, role)
    VALUE ('admin', '$2a$14$L1xC68qegONUZLHxPrDUbuVYD/niLfF7ftAkHeO3bI.5/WcIEs5vG', 'admin') # password = admin