## run_single

run single node:  

```
$ docker-compose -f docker-compose.single.yml
```

single seeds:  

```
+------------+-----------------+-------------------+--------+-------+
| token      | dc              | rack              | name   | port  |
+------------+-----------------+-------------------+--------+-------+
| 0          | asia-northeast1 | asia-northeast1-a | dyn001 | 8101  |
+------------+-----------------+-------------------+--------+-------+
```

## run_rack2

run multiple nodes:  

```
$ docker-compose -f docker-compose.run_rack2.yml
```

server seeds blow:  

```
+------------+-----------------+-------------------+--------+-------+
| token      | dc              | rack              | name   | port  |
+------------+-----------------+-------------------+--------+-------+
| 0          | asia-northeast1 | asia-northeast1-a | 1a_001 | 8101  |
| 2147483647 | asia-northeast1 | asia-northeast1-a | 1a_101 | 8102  |
| 0          | asia-northeast1 | asia-northeast1-b | 1b_001 | 8201  |
| 2147483647 | asia-northeast1 | asia-northeast1-b | 1b_101 | 8202  |
+------------+-----------------+-------------------+--------+-------+
```

describe  

```
$ ./cluster_describe.sh
{
    "dcs": [
        {
            "name": "asia-northeast1",
            "racks": [
                {
                    "name": "asia-northeast1-a",
                    "servers": [
                        {
                            "host": "dyn_1a_001",
                            "name": "192.168.16.9",
                            "port": 8101,
                            "token": 0
                        },
                        {
                            "host": "192.168.16.8",
                            "name": "192.168.16.8",
                            "port": 8101,
                            "token": 2147483647
                        }
                    ]
                },
                {
                    "name": "asia-northeast1-b",
                    "servers": [
                        {
                            "host": "192.168.16.6",
                            "name": "192.168.16.6",
                            "port": 8101,
                            "token": 0
                        },
                        {
                            "host": "192.168.16.7",
                            "name": "192.168.16.7",
                            "port": 8101,
                            "token": 2147483647
                        }
                    ]
                }
            ]
        }
    ]
}
```
