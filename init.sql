DROP TABLE IF EXISTS user;
CREATE TABLE user (
  `id`    bigint(100) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `login` varchar(128) NOT NULL UNIQUE KEY,
  `name`  varchar(128),
  `created_at`  datetime,
  `updated_at`  datetime
);
