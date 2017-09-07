# go-orderbook
Using coinbase Market data feed generating reports and redis data

# To get Only the best bid and ask into Redis (level 1)
## Run redis container
make redis

## Run server  (Right now it only does five calls at one run and adds to redis but can be expanded to run indefinitely and add to redis)
 ./orderbook
 
## To query Redis
> make client

127.0.0.1:6379> smembers BTC-USD
1) "3977253769"
2) "3977266019"
3) "3977266058"
4) "3977266120"
5) "3977266204"
6) "3977266271"
127.0.0.1:6379> HGETALL BTC-USD:3977266019
 1) "create-time"
 2) "2017-09-07T22:45:52Z"
 3) "sequence"
 4) "3977266019"
 5) "price"
 6) "4626.05"
 7) "size"
 8) "0.08"
 9) "numorders"
10) "4"

# Fetch Top 50 bids and asks into CSV files per product supported
./orderbook -d
