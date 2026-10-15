package main

import (
	"context"
	"fmt"
	"os"

	"github.com/voipbin/voipbin-go"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func main() {
	accessKey := os.Getenv("VOIPBIN_ACCESS_KEY")
	if accessKey == "" {
		fmt.Println("Set VOIPBIN_ACCESS_KEY environment variable before running")
		os.Exit(1)
	}

	client, err := voipbin.NewClient(accessKey)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	params := &voipbin_client.GetCallsParams{}
	tmp, err := client.GetCallsWithResponse(ctx, params)
	if err != nil {
		panic(err)
	}

	if tmp.JSON200 == nil {
		panic("unexpected nil response")
	}

	fmt.Printf("Next page token: %s\n", *tmp.JSON200.NextPageToken)
	for i, c := range *tmp.JSON200.Result {
		fmt.Printf("Call %d: %v\n", i, *c.Id)
	}

}
