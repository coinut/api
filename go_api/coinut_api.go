package coinut_api

import (
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "math/rand"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "bytes"
  "fmt"
)

type Credential struct {
    APIKey string
    User string
}

var Info Credential

// initializa the api with user's username and api key
func Init(user string, key string) {
    Info.User = user
    Info.APIKey = key
}

func ComputeHmac256(secret string, message string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return hex.EncodeToString(h.Sum(nil))
}

/*
    Get my balance
    Returns: your balance in a map.
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Get_balance())
    map[USDT:1.50000000 BCH:0.00000000 SGD:0.00000000 status:[OK] BTG:0.00000000 LTC:0.09678999 ZEC:0.19990000 ...]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#get-account-balance
*/
func Get_balance() map[string]interface{} {
    result, _ := Request("user_balance", "{}")
    return result
}

/*
    Get spot trading instruments information
    Args:
        pair (string): it can be any spot trading pair like "BTCUSDT"
        or "LTCBTC".
    Returns:
        if pair argument is specified, return the pair's
        information in a dict; otherwise returns all spot trading
        pairs' information
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Get_spot_instruments("LTCBTC"))
    map[inst_id:1 quote:BTC base:LTC decimal_places:5]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#get-spot-trading-instruments
*/
func Get_spot_instruments(pair string) interface{} {
    result, _ := Request("inst_list", `{"sec_type":"SPOT"}`)
    if pair != "" {
        return result["SPOT"].(map[string]interface{})[pair].([]interface{})[0]
    } else {
        return result["SPOT"]
    }
}

/*
    Get a spot trading instrument's inst_id. This inst_id is needed for
    submitting, canceling, or querying orders and for checking ticks and
    orderbooks.
    Args:
        pair (string): it can be any spot trading pair like "BTCUSDT" or "LTCBTC".
    Returns:
        return the spot trading pair's inst_id.
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Get_spot_inst_id("LTCBTC"))
    1
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#get-spot-trading-instruments
*/
func Get_spot_inst_id(pair string) int64 {
    insts := Get_spot_instruments(pair)
    return int64(insts.(map[string]interface{})["inst_id"].(float64))
}

/*
    Get a spot trading instrument's last tick.
    Args:
        inst_id (int64): the inst_id can be obtained using the
        Get_spot_inst_id or Get_spot_instruments functions.
    Returns:
        the instrument's last tick
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Get_inst_tick(1))
    map[reply:inst_tick volume:0.07546633 lowest_sell:8267.35000000 last:7797.87000000 status:[OK] ...]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#get-realtime-ticks
*/
func Get_inst_tick(inst_id int64) map[string]interface{} {
    result, _ := Request("inst_tick", fmt.Sprintf("{\"inst_id\":%d}", inst_id))
    return result
}

/*
    Get a spot trading instrument's orderbook.
    Args:
        inst_id (int64): the inst_id can be obtained using the
        Get_spot_inst_id or Get_spot_instruments functions.
    Returns:
        the trading pair's orderbook
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Get_orderbook(1))
    map[inst_id:1 sell:[map[count:1 price:8267.35 qty:0.06000000] ...] buy:[map[count:1 price:8260.35 qty:0.06000000] ...]]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#get-orderbooks-in-realtime
*/
func Get_orderbook(inst_id int64) map[string]interface{} {
    result, _ := Request("inst_order_book", fmt.Sprintf("{\"inst_id\":%d}", inst_id))
    return result
}

/*
    Get my open orders.
    Args:
        inst_id (int64): the inst_id can be obtained using the
        Get_spot_inst_id or Get_spot_instruments functions.
    Returns:
        my open orders for an instrument
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Get_open_orders(1))
    [map[order_id:454 price:0.20000000 qty:0.00110000 side:BUY client_ord_id:3.630041631e+09 inst_id:1 open_qty:0.00110000] ...]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#get-open-orders
*/
func Get_open_orders(inst_id int64) []interface{} {
    result, _ := Request("user_open_orders", fmt.Sprintf("{\"inst_id\":%d}", inst_id))
    return result["orders"].([]interface{})
}

/*
    Create a json string containing the information for opening a new order
    Args:
        inst_id (int64): the inst_id can be obtained using the
        Get_spot_inst_id or Get_spot_instruments functions.
        side (string): either 'BUY' or 'SELL'. It's case sensitive.
        qty (float64): the quantity that you want to buy or sell
        price (float64): if price is None, it's a market order; otherwise it's the price of the limit order.
        client_ord_id (int64): an id specified by the client.
    Returns:
        a map containing the information for the new order
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Create_new_order(1, "BUY", 0.0011, 0.2, -1))
    {"client_ord_id":3630041631,"inst_id":1,"price":"0.20000000","qty":"0.00110000","side":"BUY"}
*/
func Create_new_order(inst_id int64,
                      side string,
                      qty float64,
                      price float64,
                      client_ord_id int64) string {
    order := make(map[string]interface{})
    order["inst_id"] = inst_id
    order["side"] = side
    order["qty"] = fmt.Sprintf("%.8f", qty)
    if price > 0.0  {
        order["price"] = fmt.Sprintf("%.8f", price)
    }
    if client_ord_id > 0 {
        order["client_ord_id"] = client_ord_id
    } else {
        order["client_ord_id"] = rand.Int63n(4294967290) + 1

    }
    data, _ := json.Marshal(order)
    return string(data[:])
}
/*
    Submit an order to the exchange
    Args:
        inst_id (int64): the inst_id can be obtained using the get_spot_inst_id or get_spot_instruments functions.
        side (string): either 'BUY' or 'SELL'. It's case sensitive.
        qty (float64): the quantity that you want to buy or sell
        price (float64): use None to specifcy that the order is a market order; otherwise it's the price of the limit order.
        client_ord_id (int64): an id specified by the client.
    Returns:
        An order can be rejected, partially filled, or filled. Check https://github.com/coinut/api/wiki/Websocket-API#submit-an-order for the formats
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Submit_new_order(1, "BUY", 0.0011, 0.2, 2432432))
    map[qty:0.00110000 status:[OK] open_qty:0.00110000 order_id:4.5439e+06 price:0.20000000 reply:order_accepted side:BUY inst_id:1 ...]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#submit-an-order
*/
func Submit_new_order(inst_id int64,
                      side string,
                      qty float64,
                      price float64,
                      client_ord_id int64) map[string]interface{} {
    result, _ := Request("new_order", Create_new_order(inst_id, side, qty, price, client_ord_id))
    return result
}
/*
    Submit orders to the exchange
    Args:
        ords ([]string): a slice of orders in JSON string
    Returns:
        An order can be rejected, partially filled, or filled. Check https://github.com/coinut/api/wiki/Websocket-API#submit-an-order for the formats
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> orders := []string{coinut_api.Create_new_order(1, "BUY", 0.0011, 0.2, -1), coinut_api.Create_new_order(1, "BUY", 0.0015, 0.1, -1)}
    >> fmt.Println(coinut_api.Submit_new_orders())
    map[]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#submit-an-order
*/
func Submit_new_orders(ords []string) map[string]interface{} {
    data := make(map[string]interface{})
    var orders [](map[string]interface{})
    for _, ord := range ords {
        order := make(map[string]interface{})
        json.Unmarshal([]byte(ord), &order)
        orders = append(orders, order)
    }
    data["orders"] = orders
    dt, _ := json.Marshal(data)
    result, _ := Request("new_orders", string(dt[:]))
    return result
}
/*
    Cancel an order
    Args:
        inst_id (int64): the inst_id can be obtained using the get_spot_inst_id or get_spot_instruments functions.
        order_id (int64): the id of the order to be canceled
    Returns:
        Order cancel result
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Cancel_order(1, 4543903))
    map[client_ord_id:2.432432e+06 order_id:4.543903e+06 reply:cancel_order status:[OK] trans_id:7.1125678e+07 ...]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#cancel-an-order
*/
func Cancel_order(inst_id int, order_id int64) map[string]interface{} {
    result, _ := Request("cancel_order", fmt.Sprintf("{\"inst_id\":%d,\"order_id\":%d}", inst_id, order_id))
    return result
}
/*
    Cancel orders in a batch
    Args:
        inst_id (int64): the inst_id can be obtained using the
        Get_spot_inst_id or Get_spot_instruments functions.
        order_ids (a slice of int64): all the ids of the orders to be
        canceled. The maximum number of orders is 1000.
    Returns:
        cancel results
    Examples:
    >> coinut_api.Init('your username', 'your REST API Key on https://coinut.com/account/settings')
    >> fmt.Println(coinut_api.Cancel_orders(1, []int64{4543905, 4543906}))
    map[reply:cancel_orders results:[map[client_ord_id:24322 inst_id:1 order_id:4.543905e+06 status:OK] map[inst_id:1 order_id:4.543906e+06 status:INVALID_ORDER_ID client_ord_id:0]] status:[OK] trans_id:7.1126728e+07 nonce:9.47779411e+08]
    See also:
        https://github.com/coinut/api/wiki/Websocket-API#cancel-orders-in-batch-mode
*/
func Cancel_orders(inst_id int, order_ids []int64) map[string]interface{} {
    var ords [](map[string]interface{})
    for _, order_id := range order_ids {
      ords = append(ords, map[string]interface{}{"inst_id": inst_id, "order_id": order_id})
    }
    data := make(map[string]interface{})
    data["entries"] = ords
    orders, _ := json.Marshal(data)
    result, _ := Request("cancel_orders", string(orders[:]))
    return result
}

func Request(api string, content string) (map[string]interface{}, error) {
  url := "https://api.coinut.net"
  params := make(map[string]interface{})
  json.Unmarshal([]byte(content), &params)
  params["request"] = api
  params["nonce"] = rand.Int63n(1000000000) + 1
  data, _ := json.Marshal(params)
  sig := ComputeHmac256(Info.APIKey, string(data[:]))
  client := &http.Client{}
  req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
  req.Header.Add("X-User", Info.User)
  req.Header.Add("X-SIGNATURE", sig)
  resp, err := client.Do(req)
  result := make(map[string]interface{})
  if err != nil {
    return result, err
  } else {
    defer resp.Body.Close()
    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    json.Unmarshal(bodyBytes, &result)
    return result, nil
  }
}

