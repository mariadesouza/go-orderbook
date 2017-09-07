# go-orderbook
Using coinbase Market data feed generating reports and redis data

# To get Only the best bid and ask into Redis (level 1)
> make redis
> ./orderbook
## To query Redis
> make client
> 127.0.0.1:6379> HGETALL BTC-USD

# Fetch Top 50 bids and asks into CSV files per product supported
./orderbook -d
