package main

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fogfish/blueprint-serverless-golang/internal/services/restapi"
	httpd "github.com/fogfish/gouldian/v2/server/aws/apigateway"
)

func main() {
	api := restapi.NewPetShopAPI()

	handler := httpd.Serve(
		api.List(),
		api.Lookup(),
	)

	lambda.Start(func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		r.Path = strings.TrimPrefix(r.Path, "/api")
		return handler(r)
	})
}
