package bitflier

import (
    "fmt"
    "buy-xrp/app/api_interface"
    "net/http"
    "io"
    "time"
    "strconv"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "buy-xrp/config"
    "encoding/json"
    "strings"
)

type Order struct {
	ProductCode    string  `json:"product_code"`
	ChildOrderType string  `json:"child_order_type"`
	Side           string  `json:"side"`
	Price          float64 `json:"price"`
	Size           float64 `json:"size"`
	MinuteToExpire int     `json:"minute_to_expire"`
	TimeInForce    string  `json:"time_in_force"`
}

type Apis struct {
    product_code string
    value float64
}

func NewBitFlier(code string) api_interface.Api{
    return &Apis{
        product_code:code,

    }
}

func setPrivateHeader(req *http.Request){
    req.Header.Add("ACCESS-KEY",config.Config.ApiKey)
    t := time.Now().Unix()
    ts := strconv.FormatInt(t, 10)
    req.Header.Add("ACCESS-TIMESTAMP",ts)
    req_body, _ := req.GetBody()
    body_byte, _ := io.ReadAll(req_body)
    var hira =ts+req.Method+req.URL.RequestURI()+string(body_byte)
    mac := hmac.New(sha256.New, []byte(config.Config.ApiSecret))
    mac.Write([]byte(hira))
    req.Header.Add("ACCESS-SIGN",hex.EncodeToString(mac.Sum(nil)))
    req.Header.Add("Content-Type","application/json")
}

func (api *Apis) GetTicker() api_interface.Ticker{
    fmt.Println("GetTicker")
    resp, _ := http.Get(config.BasePath+config.TickerUrl+"?product_code="+api.product_code)
    defer resp.Body.Close()
    body,_ := io.ReadAll(resp.Body)
    var decoded api_interface.Ticker
    json.Unmarshal(body, &decoded)
    return decoded
}

func (api *Apis) GetBalance() {
    fmt.Println("GetBalance")
    client := &http.Client{CheckRedirect: nil}
    req_body := strings.NewReader("")
    req, _ := http.NewRequest("GET", config.BasePath+config.MyBalance, req_body)
    setPrivateHeader(req)
    res, _ := client.Do(req)
    defer res.Body.Close()
    body,_ := io.ReadAll(res.Body)
    fmt.Println(string(body))
}

func (api *Apis) AllOrderCancel(){
    fmt.Println("AllOrderCancel")
    client := &http.Client{CheckRedirect: nil}

    req_body := strings.NewReader(`{"product_code":"`+api.product_code+`"}`)
    req, _ := http.NewRequest("POST", config.BasePath+config.MyAllOrderCancel, req_body)
    setPrivateHeader(req)

    res, _ := client.Do(req)
    defer res.Body.Close()

    body,_ := io.ReadAll(res.Body)
    fmt.Println(string(body))
}

func (api *Apis) SetValue(value float64){
    api.value = value
}

func (api *Apis) Order() string {
    fmt.Println("Order")
    client := &http.Client{CheckRedirect: nil}

    var order Order
    order.ProductCode = api.product_code
    order.ChildOrderType = "MARKET"
    order.Side = "BUY"
    order.Size = api.value
    order_json, _ := json.Marshal(order)

    req_body := strings.NewReader(string(order_json))
    req, _ := http.NewRequest("POST", config.BasePath+config.MySendOrder, req_body)
    setPrivateHeader(req)

    res, _ := client.Do(req)
    defer res.Body.Close()

    body,_ := io.ReadAll(res.Body)
    return string(body)
}

