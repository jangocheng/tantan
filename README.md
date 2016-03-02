# tantan
# work envieonment
  go1.6 + PostgreSQL 9.4.4 + redis 3.0.3
# Usage

	cd $GOPATH/src/github.com
	git clone https://github.com/lxbgit/tantan.git
	cd tantan
	go build
    

init && start postgresql:
  
  - ./bin/pg_init.py
  
start:

  - ./bin/pg_start.py
  - ./bin/redis_start.py
  - ./tantan --debug
 
stop:
  
  - ./bin/pg_stop.py
  - ./bin/redis_stop.py
   
the configure file is  conf/service.json

default http port:8080

the database sql file is conf/table.sql
    
