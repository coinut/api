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

type CoinutClient struct {
    APIKey string
    User string
}

// initialize the api with user's username and api key
func NewClient(user string, key string) *CoinutClient {
    c := &CoinutClient{APIKey: key, User: user}
    return c
}

//    Get my balance
//    Returns: your balance in a map.
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.GetBalance()
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: map[USDT:1.50000000 BCH:0.00000000 SGD:0.00000000 status:[OK] BTG:0.00000000 LTC:0.09678999 ZEC:0.19990000 ...]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#get-account-balance
func (client *CoinutClient) GetBalance() (map[string]interface{}, error) {
    return client.Request("user_balance", "{}")
}

//    Get spot trading instruments information
//    Args:
//        pair (string): it can be any spot trading pair like "BTCUSDT" or "LTCBTC".
//
//    Returns:
//        if pair argument is specified, return the pair's
//        information in a dict; otherwise returns all spot trading
//        pairs' information
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.GetSpotInstruments("LTCBTC")
//        if err == nil {
//            fmt.Println(result)
//        }
//        map[inst_id:1 quote:BTC base:LTC decimal_places:5]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#get-spot-trading-instruments
func (client *CoinutClient) GetSpotInstruments(pair string) (interface{}, error) {
    result, err := client.Request("inst_list", `{"sec_type":"SPOT"}`)
    if err != nil {
        return result, err
    }
    if pair != "" {
        return result["SPOT"].(map[string]interface{})[pair].([]interface{})[0], nil
    } else {
        return result["SPOT"], nil
    }
}

//    Get a spot trading instrument's inst_id. This inst_id is needed for
//    submitting, canceling, or querying orders and for checking ticks and
//    orderbooks.
//    Args:
//        pair (string): it can be any spot trading pair like "BTCUSDT" or "LTCBTC".
//
//    Returns:
//        return the spot trading pair's inst_id.
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings"
//        result, err := client.GetSpotInstId("LTCBTC")
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: 1
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#get-spot-trading-instruments
func (client *CoinutClient) GetSpotInstId(pair string) (uint32, error) {
    insts, err := client.GetSpotInstruments(pair)
    if err != nil {
        return 0, err
    }
    return uint32(insts.(map[string]interface{})["inst_id"].(float64)), nil
}

//    Get a spot trading instrument's last tick.
//    Args:
//        inst_id (uint32): the inst_id can be obtained using the GetSpotInstId or GetSpotInstruments functions.
//
//    Returns:
//        the instrument's last tick
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.GetInstTick(1)
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: map[reply:inst_tick volume:0.07546633 lowest_sell:8267.35000000 last:7797.87000000 status:[OK] ...]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#get-realtime-ticks
func (client *CoinutClient) GetInstTick(inst_id uint32) (map[string]interface{}, error) {
    return client.Request("inst_tick", fmt.Sprintf("{\"inst_id\":%d}", inst_id))
}

//    Get a spot trading instrument's orderbook.
//    Args:
//        inst_id (uint32): the inst_id can be obtained using the GetSpotInstId or GetSpotInstruments functions.
//
//    Returns:
//        the trading pair's orderbook
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.GetOrderbook(1)
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: map[inst_id:1 sell:[map[count:1 price:8267.35 qty:0.06000000] ...] buy:[map[count:1 price:8260.35 qty:0.06000000] ...]]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#get-orderbooks-in-realtime
func (client *CoinutClient) GetOrderbook(inst_id uint32) (map[string]interface{}, error) {
    return client.Request("inst_order_book", fmt.Sprintf("{\"inst_id\":%d}", inst_id))
}

//    Get my open orders.
//    Args:
//        inst_id (uint32): the inst_id can be obtained using the GetSpotInstId or GetSpotInstruments functions.
//
//    Returns:
//        my open orders for an instrument
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.GetOpenOrders(1)
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: [map[order_id:454 price:0.20000000 qty:0.00110000 side:BUY client_ord_id:3.630041631e+09 inst_id:1 open_qty:0.00110000] ...]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#get-open-orders
func (client *CoinutClient) GetOpenOrders(inst_id uint32) ([]interface{}, error) {
    result, err := client.Request("user_open_orders", fmt.Sprintf("{\"inst_id\":%d}", inst_id))
    if err != nil {
        return make([]interface{}, 0), err
    } else {
        return result["orders"].([]interface{}), nil
    }
}

//    Create a json string containing the information for opening a new order
//    Args:
//        inst_id (uint32): the inst_id can be obtained using the GetSpotInstId or GetSpotInstruments functions.
//        side (string): either 'BUY' or 'SELL'. It's case sensitive.
//        qty (float64): the quantity that you want to buy or sell
//        price (float64): if price is None, it's a market order; otherwise it's the price of the limit order.
//        client_ord_id (uint32): an id specified by the client.
//
//    Returns:
//        a map containing the information for the new order
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.CreateNewOrder(1, "BUY", 0.0011, 0.2, 0)
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: {"client_ord_id":3630041631,"inst_id":1,"price":"0.20000000","qty":"0.00110000","side":"BUY"}
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#ceate-new-order
func (client *CoinutClient) CreateNewOrder(inst_id uint32,
                      side string,
                      qty float64,
                      price float64,
                      client_ord_id uint32) (string, error) {
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
        order["client_ord_id"] = rand.Uint32()
    }
    data, _ := json.Marshal(order)
    return string(data[:]), nil
}

//    Submit an order to the exchange
//    Args:
//        inst_id (uint32): the inst_id can be obtained using the get_spot_inst_id or get_spot_instruments functions.
//        side (string): either 'BUY' or 'SELL'. It's case sensitive.
//        qty (float64): the quantity that you want to buy or sell
//        price (float64): use None to specifcy that the order is a market order; otherwise it's the price of the limit order.
//        client_ord_id (uint32): an id specified by the client.
//
//    Returns:
//        An order can be rejected, partially filled, or filled. Check https://github.com/coinut/api/wiki/Websocket-API#submit-an-order for the formats
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.SubmitNewOrder(1, "BUY", 0.0011, 0.2, 2432432)
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: map[qty:0.00110000 status:[OK] open_qty:0.00110000 order_id:4.5439e+06 price:0.20000000 reply:order_accepted side:BUY inst_id:1 ...]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#submit-an-order
func (client *CoinutClient) SubmitNewOrder(inst_id uint32,
                      side string,
                      qty float64,
                      price float64,
                      client_ord_id uint32) (map[string]interface{}, error) {
    order, err := client.CreateNewOrder(inst_id, side, qty, price, client_ord_id)
    if err != nil {
        return make(map[string]interface{}), err
    }
    return client.Request("new_order", order)
}

//    Submit orders to the exchange
//    Args:
//        ords ([]string): a slice of orders in JSON string
//
//    Returns:
//        An order can be rejected, partially filled, or filled. Check https://github.com/coinut/api/wiki/Websocket-API#submit-an-order for the formats
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        order1, _ := client.CreateNewOrder(1, "BUY", 0.0011, 0.2, 0)
//        order2, _ := client.CreateNewOrder(1, "SELL", 0.02, 0.2, 0)
//        orders := []string{order1, order2}
//        result, err := client.SubmitNewOrders(orders)
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: map[orders:[{"client_ord_id":2596996162,"inst_id":1,"nonce":2854263694,"open_qty":"0.00110000","order_id":4543906,"price":"0.20000000","qty":"0.00110000","reply":"order_accepted","side":"BUY","status":["OK"],"timestamp":1530637392621081,"trans_id":71160389}, {"client_ord_id":4039455774,"inst_id":1,"nonce":2854263694,"open_qty":"0.00150000","order_id":4543907,"price":"0.10000000","qty":"0.00150000","reply":"order_accepted","side":"BUY","status":["OK"],"timestamp":1530637392621081,"trans_id":71160391}]]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#submit-orders-in-batch-mode
func (client *CoinutClient) SubmitNewOrders(ords []string) ([]interface{}, error) {
    data := make(map[string]interface{})
    var orders [](map[string]interface{})
    for _, ord := range ords {
        order := make(map[string]interface{})
        json.Unmarshal([]byte(ord), &order)
        orders = append(orders, order)
    }
    data["orders"] = orders
    dt, _ := json.Marshal(data)
    result, err := client.Request("new_orders", string(dt[:]))
    if err == nil {
        return result["orders"].([]interface{}), nil
    } else {
        return make([]interface{}, 0), err
    }
}

//    Cancel an order
//    Args:
//        inst_id (uint32): the inst_id can be obtained using the get_spot_inst_id or get_spot_instruments functions.
//        order_id (uint32): the id of the order to be canceled
//
//    Returns:
//        Order cancel result
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.CancelOrder(1, 4543903)
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: map[client_ord_id:2.432432e+06 order_id:4.543903e+06 reply:cancel_order status:[OK] trans_id:7.1125678e+07 ...]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#cancel-an-order
func (client *CoinutClient) CancelOrder(inst_id int, order_id uint32) (map[string]interface{}, error) {
    return client.Request("cancel_order", fmt.Sprintf("{\"inst_id\":%d,\"order_id\":%d}", inst_id, order_id))
}

//    Cancel orders in a batch
//    Args:
//        inst_id (uint32): the inst_id can be obtained using the Get_spot_inst_id or Get_spot_instruments functions.
//        order_ids (a slice of uint32): all the ids of the orders to be canceled. The maximum number of orders is 1000.
//
//    Returns:
//        cancel results
//
//    Examples:
//        import "github.com/coinut/api/go_api"
//        client := coinut_api.NewClient("your username", "your REST API Key on https://coinut.com/account/settings")
//        result, err := client.CancelOrders(1, []uint32{4543905, 4543906})
//        if err == nil {
//            fmt.Println(result)
//        }
//        output: map[reply:cancel_orders results:[map[client_ord_id:24322 inst_id:1 order_id:4.543905e+06 status:OK] map[inst_id:1 order_id:4.543906e+06 status:INVALID_ORDER_ID client_ord_id:0]] status:[OK] trans_id:7.1126728e+07 nonce:9.47779411e+08]
//
//    See also:
//        https://github.com/coinut/api/wiki/Websocket-API#cancel-orders-in-batch-mode
func (client *CoinutClient) CancelOrders(inst_id int, order_ids []uint32) (map[string]interface{}, error) {
    var ords [](map[string]interface{})
    for _, order_id := range order_ids {
      ords = append(ords, map[string]interface{}{"inst_id": inst_id, "order_id": order_id})
    }
    data := make(map[string]interface{})
    data["entries"] = ords
    orders, _ := json.Marshal(data)
    return client.Request("cancel_orders", string(orders[:]))
}

func ComputeHmac256(secret string, message string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return hex.EncodeToString(h.Sum(nil))
}

func (client *CoinutClient) Request(api string, content string) (map[string]interface{}, error) {
    url := "https://api.coinut.com"
    params := make(map[string]interface{})
    json.Unmarshal([]byte(content), &params)
    params["request"] = api
    params["nonce"] = rand.Int63n(4294967200) + 1
    data, _ := json.Marshal(params)
    sig := ComputeHmac256(client.APIKey, string(data[:]))
    cli := &http.Client{}
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
    req.Header.Add("X-User", client.User)
    req.Header.Add("X-SIGNATURE", sig)
    resp, err := cli.Do(req)
    var result interface{}
    if err != nil {
        return make(map[string]interface{}), err
    } else {
        defer resp.Body.Close()
        bodyBytes, _ := ioutil.ReadAll(resp.Body)
        err := json.Unmarshal(bodyBytes, &result)
        if err != nil {
            return make(map[string]interface{}), err
        }
        r, ok := result.(map[string]interface{})
        if ok {
            return r, nil
        } else {
            return map[string]interface{}{"orders": result}, nil
        }
    }
}

