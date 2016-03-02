# tantan
# work envieonment
  go1.6 + PostgreSQL 9.4.4 + redis 3.0.3
## Usage

    cd tantan
    
  init && start postgresql:
    ./bin/pg_init.py
    
  start:
    (1)./bin/pg_start.py
    (2)./bin/redis_start.py
    (3)go build
    (4)./tantan --debug
   stop:
    (1)./bin/pg_stop.py
    (2)./bin/redis_stop.py
   
   the configure file is  conf/service.json
   the database sql file is conf/table.sql
    