package main

import "context"

type app interface {
	Run(ctx context.Context)
}
