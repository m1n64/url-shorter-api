package chromium

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"sync"
)

type ChromeAllocator struct {
	AllocatorCtx context.Context
	Cancel       context.CancelFunc
	once         sync.Once
}

func NewChromeAllocator() *ChromeAllocator {
	return &ChromeAllocator{}
}

func (ca *ChromeAllocator) Init() {
	ca.once.Do(func() {
		ctx, cancel := chromedp.NewExecAllocator(context.Background(),
			append(chromedp.DefaultExecAllocatorOptions[:],
				chromedp.Flag("headless", true),
				chromedp.Flag("disable-gpu", true),
				chromedp.Flag("disable-extensions", true),
				chromedp.Flag("no-sandbox", true),
				chromedp.Flag("disable-setuid-sandbox", true),
			)...,
		)

		ca.AllocatorCtx = ctx
		ca.Cancel = cancel

		log.Println("Chrome Allocator has been initialized.")
	})
}

func (ca *ChromeAllocator) Close() {
	if ca.Cancel != nil {
		ca.Cancel()
		log.Println("Chrome Allocator has been closed.")
	}
}
