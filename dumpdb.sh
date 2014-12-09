#!/bin/bash

mysqldump --user=root --password=root fabsmith > fabsmith.sql

# GRANT ALL PRIVILEGES ON fabsmith.* To 'fabsmith'@'localhost' IDENTIFIED BY 'fabsmith';