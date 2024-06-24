#!/usr/bin/env bash
options='-config=./internal/database/migrations/sql-migrate.yml'
sql-migrate $@ $options