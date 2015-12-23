# Hello SQL

This is a test SQL REST API test tool, used for SQL benchmark

## SQL Schema

```
CREATE TABLE entries (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `value` VARCHAR(255),
    PRIMARY KEY(`id`)
) ENGINE=InnoDB CHARACTER SET utf8 COLLATE utf8_general_ci;

```

## Configuration

Use environment:

```
SQL_DRIVER=mysql
SQL_DSN=user:password@tcp(127.0.0.1:3306)/hello_sql
```
