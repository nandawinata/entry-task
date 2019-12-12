#!/bin/bash
set -e
service mysql start
mysql < /scripts/sql/init_table.sql
service mysql stop