#!/usr/bin/bash

# Registers a user and authenticats it

curl -XPOST http://138.124.107.242:80/api/register -F login=testtest -F pswd=Abcde_12345 -F token=abcde12345
curl -XPOST http://138.124.107.242:80/api/auth -F login=testtest -F pswd=Abcde_12345
