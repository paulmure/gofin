-- US Stock exchange uses Eastern Time
SET GLOBAL time_zone = '-5:00';

CREATE DATABASE us_stocks;
USE us_stocks;

CREATE TABLE ticker (
	id INT NOT NULL AUTO_INCREMENT,
    symbol VARCHAR(6) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE intraday_by_minute (
	id INT NOT NULL AUTO_INCREMENT,
    ticker_id INT NOT NULL,
    quote_timestamp DATETIME NOT NULL,
    bid DECIMAL(12, 2) NOT NULL,
    ask DECIMAL(12, 2) NOT NULL,
    bid_size INT NOT NULL,
    ask_size INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (ticker_id) REFERENCES ticker(id)
);