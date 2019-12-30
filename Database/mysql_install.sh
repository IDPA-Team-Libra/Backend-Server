#!/bin/bash
service mysql start
mysql -u root --password= < /mysql/setup.sql
mysqladmin -u root password root_password_123
service mysql stop
