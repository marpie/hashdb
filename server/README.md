hashdb server
=============

This is a sample http server.

It understands three HTTP Requests:
  - /md5?q=AAAAAAA      -> (GET) searches the MD5 hash in the database
  - /sha1?q=BBBBBB      -> (GET) searches the SHA1 hash in the database
  - /new?p=newPassword  -> (POST) adds the password to the database

Before the server can be started for the first time the following 
directories must exist:
  - db
    - db/md5
    - db/sha1

Use `mkdir db db/md5 db/sha1` in the server directory to create them.

