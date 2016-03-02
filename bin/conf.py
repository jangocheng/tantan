#!/usr/bin/env python
#-*- coding:utf-8 -*-

import os
import json

data = open(os.path.realpath("./conf/service.json")).read().strip()
config_info = json.loads(data)


def get_db_config():
    return config_info["database"]


def get_redis_config():
    return config_info["redis"]
