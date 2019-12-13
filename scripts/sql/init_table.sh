#!/bin/bash
set -e
service mysql start
mysql < init_table.sql
service mysql stop