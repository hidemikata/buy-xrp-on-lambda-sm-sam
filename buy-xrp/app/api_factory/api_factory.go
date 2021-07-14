package api_factory

type ApiExchanger struct{
    ExchangerApi func(a string) api_interface.Api
}



