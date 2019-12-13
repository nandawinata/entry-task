#!/bin/bash
set -e
service mysql start
mysql < /root/sql/init_table.sql
service mysql stop