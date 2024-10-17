package rag

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/vogtp/rag/pkg/scraper"
)

func WebScrapRag(ctx context.Context, url string) (llms.Model, error) {
	scrap,err:=scraper.New()
	if err != nil {
		return nil, err
	}
	rsp, err :=scrap.Call(ctx, url )
	if err != nil {
		return nil, err
	}
	fmt.Println(rsp)
	return nil, nil
}
