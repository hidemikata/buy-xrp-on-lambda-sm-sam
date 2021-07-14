package config
import(
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
    )

const BasePath="https://api.bitflyer.com"

//public api
const TickerUrl="/v1/getticker"

//private api
//資産情報
const MyBalance = "/v1/me/getbalance"

const MyAllOrderCancel = "/v1/me/cancelallchildorders"
const MySendOrder = "/v1/me/sendchildorder"
const BuyValue = "1000"



type ConfigList struct{
    ApiKey string
    ApiSecret string
    BuyValue string//yen
}
var Config ConfigList

func init(){
    Config.ApiKey, _ = getParameter("bitflyer-api-key")
    Config.ApiSecret, _ = getParameter("bitflyer-secret-key")
}

func getParameter(key string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ssm.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

	params := &ssm.GetParameterInput{
		Name:			aws.String(key),
		WithDecryption: aws.Bool(true),
	}

	res, err := svc.GetParameter(params)
	if err != nil {
		return "", err
	}

	return *res.Parameter.Value, nil
}
