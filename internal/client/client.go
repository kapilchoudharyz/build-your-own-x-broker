package client

import "context"

type Client struct {
	messages chan string
	ctx      context.Context
	connectionID string
}