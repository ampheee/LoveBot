#!/bin/bash

apt update -y
apt download -y postgresql
apt download -y golang
psql -U postgres < init.sql
