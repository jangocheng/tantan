#!/usr/bin/env python
#-*- coding:utf-8 -*-

import os
import conf


def main():
    db_config = conf.get_db_config()

    for name, config in db_config.items():
        port = config["port"]
        base_dir = config['dir']
        base_dir = os.path.abspath(base_dir)
        data_dir = os.path.join(base_dir, "pg_data/" + name)
        print data_dir
        os.system('pg_ctl start -D %s -o "-p %s" -l %s/postgresql.log' % (data_dir, port, data_dir))


if __name__ == "__main__":
    main()
