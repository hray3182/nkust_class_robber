package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func LoadCookies() {
	// set headless false
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	// 1. navigate to the page
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// navigate to the page
	if err := chromedp.Run(ctx,
		chromedp.Navigate(`https://webap.nkust.edu.tw`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("登入成功後，請按Enter鍵繼續...")
			var input string
			fmt.Scanln(&input)
			return nil
		}),
		chromedp.Navigate(`https://aais3.nkust.edu.tw/selcrs_std/Home/About`),

		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(2 * time.Second)
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				log.Fatal("get cookies error: ", err)
			}
			for _, cookie := range cookies {
				log.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
			}
			// save cookie to file
			file, err := os.Create("cookie.json")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			enc := json.NewEncoder(file)
			if err := enc.Encode(cookies); err != nil {
				log.Fatal(err)
			}

			return nil
		}),
	); err != nil {
		log.Fatal(err)
	}
}
