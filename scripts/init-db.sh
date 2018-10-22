#!/usr/bin/env bash

psql -d postgres -f./pkg/db/create/1_create_db.sql
