package trades

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Bid struct {
	Price     string
	Size      string
	NumOrders int
}

func (b *Bid) csvWriter(writer *csv.Writer) {
	var record []string
	record = append(record, b.Price)
	record = append(record, b.Size)
	record = append(record, fmt.Sprintf("%d", b.NumOrders))
	writer.Write(record)
}

//Custom UnmarshalJSON
func (b *Bid) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	b.Price, _ = v[0].(string)
	b.Size, _ = v[1].(string)
	b.NumOrders = int(v[2].(float64))

	return nil
}

type OrderBook struct {
	Sequence int64 `json:"sequence"`
	Bids     []Bid `json:"bids"`
	Asks     []Bid `json:"asks"`
}

// GetOrderBook
func GetOrderBook(level, id string) (*OrderBook, error) {

	url := "https://api.gdax.com/products/" + id + "/book"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("level", level)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orders OrderBook

	err = json.Unmarshal(body, &orders)
	if err != nil {
		return nil, err
	}

	return &orders, nil
}

func (o *OrderBook) WriteToCSV(path string) error {
	os.MkdirAll(path, 0755)
	prefix := "orderBook-" + fmt.Sprintf("%d", o.Sequence)
	fullpath := filepath.Join(path, prefix+"-BIDS.csv")
	f1, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	writerBids := csv.NewWriter(f1)

	fullpath2 := filepath.Join(path, prefix+"-ASKS.csv")
	f2, err := os.Create(fullpath2)
	if err != nil {
		return err
	}
	writerAsks := csv.NewWriter(f2)

	writerBids.Write([]string{"Price", "Size", "NumOrders"})
	for _, bids := range o.Bids {
		bids.csvWriter(writerBids)
	}
	writerBids.Flush()

	writerAsks.Write([]string{"Price", "Size", "NumOrders"})
	for _, asks := range o.Bids {
		asks.csvWriter(writerAsks)
	}
	writerAsks.Flush()

	return nil
}

func (o *OrderBook) RecordTopBidinRedis(productID string) error {

	redisConn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer redisConn.Close()

	//createTime := time.Now().UTC().Format(time.RFC3339)
	//day := time.Now().UTC().Format("2006-01-02")
	var maxBid *Bid
	for _, b := range o.Bids {
		if maxBid == nil || b.Price > maxBid.Price {
			maxBid = &b
		}
	}

	//redis.ScanStruct(src, dest)
	if _, err := redisConn.Do("HMSET", productID,
		"create-time", time.Now().UTC().Format(time.RFC3339),
		"sequence", strconv.FormatInt(o.Sequence, 10),
		"price", maxBid.Price,
		"size", maxBid.Size,
		"numorders", strconv.Itoa(maxBid.NumOrders)); err != nil {
		return err
	}
	fmt.Println("Recorded in Redis")
	//127.0.0.1:6378> HGETALL LTC-EUR
	return nil
}
