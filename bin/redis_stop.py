#!/usr/bin/env python

import os

import conf

redis_config = conf.get_redis_config()

for name, config in redis_config.items():
    host = config["host"]
    port = config["port"]
    print "stop %s port: %s" % (name, port)
    os.system('redis-cli -h %s -p %d shutdown' % (host, port))