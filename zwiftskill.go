package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler() (string, error) {
	return fmt.Sprintf("Hello World"), nil
}

func main() {
	lambda.Start(Handler)
}
