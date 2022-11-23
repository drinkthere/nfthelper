#!/bin/bash
# create database and user and grant privileges to user
create database nfthelper;
create user nfthelper identified by '123456';
grant all privileges on nfthelper.* to nfthelper@'%';

# create table

