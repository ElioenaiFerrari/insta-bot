package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ElioenaiFerrari/insta-bot/internal"
	"github.com/chromedp/chromedp"
)

func WithRetries(retries int32, attempt int32) {
	backoff := time.Second * 1 * time.Duration(1<<uint(attempt))
	for range time.Tick(backoff) {
		if attempt >= retries {
			break
		}
		// Configuração do Chromedp
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", internal.Env.Env == "prd"),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
		)
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		// Cria contexto do Chromedp
		ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		defer cancel()

		// Login no Instagram
		err := chromedp.Run(ctx,
			chromedp.Navigate(`https://www.instagram.com/accounts/login/`),
			chromedp.WaitVisible(`input[name="username"]`),
			chromedp.SendKeys(`input[name="username"]`, internal.Env.User.Email),
			chromedp.SendKeys(`input[name="password"]`, internal.Env.User.Password),
			chromedp.Click(`button[type="submit"]`),
			chromedp.WaitVisible("._ac8f"),
			chromedp.Click("._ac8f"),
			chromedp.WaitVisible("button._a9--"),
			chromedp.Click("button._a9--"),
		)
		if err != nil {
			fmt.Println("Erro ao fazer login:", err)
			WithRetries(retries, attempt+1)
			return
		}

		fmt.Println("Login realizado com sucesso")

		// Abrir post e curtir
		err = chromedp.Run(ctx,
			chromedp.Navigate(internal.Env.Post.URL),
			chromedp.WaitVisible(`._aamw button[class="_abl-"]`),
			chromedp.Click(`._aamw button[class="_abl-"]`),
		)
		if err != nil {
			fmt.Println("Erro ao curtir post:", err)
			WithRetries(retries, attempt+1)
			return
		}

		fmt.Println("Post curtido com sucesso")

		for i := 0; i < internal.Env.Post.CommentTimes; i++ {
			// Deixar comentário
			err = chromedp.Run(ctx,
				chromedp.WaitVisible(`._aamx button[class="_abl-"]`),
				chromedp.Click(`._aamx button[class="_abl-"]`),
				chromedp.SetValue(`textarea`, ""),
				chromedp.SendKeys(`textarea`, internal.Env.Post.Comment),
				chromedp.KeyEvent("\r"),
			)
			if err != nil {
				fmt.Println("Erro ao deixar comentário:", err)
				WithRetries(retries, attempt+1)
				return
			}
			fmt.Printf("Comentário %d deixado com sucesso\n", i+1)
			time.Sleep(time.Minute)
		}
		break
	}
}

func main() {
	WithRetries(3, 1)
}
