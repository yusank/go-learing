# redis 压测性能

## 环境

```shell
# 宿主机
MacBook Pro (13-inch, 2018, Four Thunderbolt 3 Ports)
CPU: 2.3 GHz 四核Intel Core i5
内存:16 GB 2133 MHz LPDDR3
# redis
Redis 6.0.10 (00000000/0) 64 bit
```

## 压测性能

> 50并发 20w请求 key 长度 10

10 byte

```shell
➜redis-benchmark -h 127.0.0.1 -p 6379 -c 50 -n 20000 -d 10 -t get,set -r 10 --csv
"SET","93457.95"
"GET","96153.84"
```

20 byte

```shell
➜redis-benchmark -h 127.0.0.1 -p 6379 -c 50 -n 20000 -d 20 -t get,set -r 10 --csv
"SET","92165.90"
"GET","96153.84"
```

50 byte

```shell
➜redis-benchmark -h 127.0.0.1 -p 6379 -c 50 -n 20000 -d 50 -t get,set -r 10 --csv
"SET","91324.20"
"GET","94339.62"
```

100 byte

```shell
➜redis-benchmark -h 127.0.0.1 -p 6379 -c 50 -n 20000 -d 100 -t get,set -r 10 --csv
"SET","87719.30"
"GET","93339.62"
```

200 byte

```shell
➜redis-benchmark -h 127.0.0.1 -p 6379 -c 50 -n 20000 -d 200 -t get,set -r 10 --csv
"SET","89285.71"
"GET","88888.89"
```

1k byte

```shell
➜redis-benchmark -h 127.0.0.1 -p 6379 -c 50 -n 20000 -d 1000 -t get,set -r 10 --csv
"SET","88495.58"
"GET","91743.12"
```

5k byte

```shell
➜redis-benchmark -h 127.0.0.1 -p 6379 -c 50 -n 20000 -d 5000 -t get,set -r 10 --csv
"SET","84745.77"
"GET","89686.09"
```

表格(request per seconds)：

| op | 10 | 20 | 50 | 100 | 200 | 1k | 5K |
| ---- | :----: | :----: | :----: | :----: | :----: | :----: | :----: |
|SET|93457.95 |92165.90 |91324.20 |87719.30 |89285.71 |88495.58 |84745.77|
|GET|96153.84 |96153.84 |94339.62 |94339.62 |88888.89|91743.12 | 89686.09 |


## 测试key大小

|value size| insert num |before insert | after insert | avg used_mem| key avg used_mem |
| :----: | :----: | :----: | :----: | :----: |:----: |
|10| 10w |1066048|11714624|106.48576| 96.49|
|20 | 10w | 1066272 |13314848 | 122.48576|102.49 |
|50| 10w| 1066464 |16515040 | 154.48576 | 104.49|
|100 | 10w |1066656|21315232 | 202.48576| 102.49|
|200 | 10w |1066848 |30915424 |298.48576 | 98.49|
|1k | 10w |1067040 |110915616 | 1098.48576| 98.49|
|5k | 10w |1067232 | 522115808| 5210.48576| 210.49|
