package main

import (
	"fmt"


	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

    "buy-xrp/app/bitflier"
    "buy-xrp/config"
    "strconv"
)
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    api := bitflier.NewBitFlier("XRP_JPY")
    ticker := api.GetTicker()
    fmt.Println("best_ask", ticker.BestAsk)
    buyValueYen, _ := strconv.ParseFloat(config.BuyValue, 64)
    buyValue :=  buyValueYen / ticker.BestAsk
    fmt.Println(buyValue)
    api.SetValue(float64(int(buyValue)))
    fmt.Println(float64(int(buyValue)))

    api.AllOrderCancel()
    res := api.Order()

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello world %s", res),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

