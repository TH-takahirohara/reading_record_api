CREATE TABLE IF NOT EXISTS readings (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  book_name VARCHAR(500) NOT NULL,
  book_author VARCHAR(500) NOT NULL,
  total_page_count INT NOT NULL,
  current_page INT NOT NULL DEFAULT 0,
  finished BOOLEAN NOT NULL DEFAULT FALSE,
  memo VARCHAR(10000) NOT NULL DEFAULT '',
  user_id BIGINT UNSIGNED NOT NULL REFERENCES users ON DELETE CASCADE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  version BIGINT NOT NULL DEFAULT 1,
  PRIMARY KEY (id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;