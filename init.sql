DROP TABLE IF EXISTS user_photo;
DROP TABLE IF EXISTS user;

CREATE TABLE user (
  `id`    bigint(100) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `login` varchar(128) NOT NULL UNIQUE KEY,
  `name`  varchar(128) NOT NULL,
  `created_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;;

CREATE TABLE user_photo (
  `user_id` bigint(100) NOT NULL UNIQUE,
  `uuid` varchar(32) NOT NULL,

  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES user (id)
    ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB;
