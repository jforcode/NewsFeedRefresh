CREATE TABLE news_api_flags (
  _id INTEGER PRIMARY KEY AUTO_INCREMENT,
  key VARCHAR(255),
  value VARCHAR(255),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  status VARCHAR(255) DEFAULT 'Active'
) CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

INSERT INTO news_api_flags (key, value)
VALUES ('remaining_requests', '1000');

CREATE TABLE api_sources (
  _id INTEGER PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100),
  domain_url VARCHAR(255),
  api_home_url VARCHAR(255),
  api_url VARCHAR(255),
  attribution_name VARCHAR(255),
  attribution_label VARCHAR(255),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  status VARCHAR(255) DEFAULT 'Active'
) CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

INSERT INTO api_sources (name, domain_url, api_home_url, api_url, attribution_name, attribution_label)
VALUES ('news_api', 'https://newsapi.org', 'https://newsapi.org', 'https://newsapi.org/v2', 'News API', 'Powered by News API');

CREATE TABLE sources (
  _id INTEGER PRIMARY KEY AUTO_INCREMENT,
  api_source_name VARCHAR(100) COMMENT 'refers api_source(name)',
  s_id VARCHAR(255),
  name VARCHAR(255),
  description TEXT,
  url VARCHAR(255),
  category VARCHAR(100),
  language VARCHAR(10),
  country VARCHAR(10),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  status VARCHAR(255) DEFAULT 'Active'
) CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE articles (
  _id INTEGER PRIMARY KEY AUTO_INCREMENT,
  api_source_name VARCHAR(100) COMMENT 'refers api_source(name)',
  source_id VARCHAR(255) COMMENT 'refers source(s_id)',
  source_name VARCHAR(255) COMMENT 'refers source(name)',
  author VARCHAR(255),
  title TEXT,
  description TEXT,
  url VARCHAR(255),
  url_to_image VARCHAR(255),
  published_at DATETIME,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  status VARCHAR(255) DEFAULT 'Active'
) CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
