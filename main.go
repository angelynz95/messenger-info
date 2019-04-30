package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/paked/messenger"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load(os.Getenv("GOPATH") + "/src/github.com/angelynz95/djook/.env")

	verifyToken := flag.String("verify-token", os.Getenv("FB_BOT_VERIFY_TOKEN"), "The token used to verify facebook (required)")
	verify := flag.Bool("should-verify", false, "Whether or not the app should verify itself")
	pageToken := flag.String("page-token", os.Getenv("FB_PAGE_ACCESS_TOKEN"), "The token that is used to verify the page on facebook")
	appSecret := flag.String("app-secret", os.Getenv("FB_APP_SECRET"), "The app secret from the facebook developer portal (required)")
	flag.Parse()

	if *verifyToken == "" || *appSecret == "" || *pageToken == "" {
		fmt.Println("missing arguments")
		fmt.Println()
		flag.Usage()

		os.Exit(-1)
	}

	// Create a new messenger client
	client := messenger.New(messenger.Options{
		Verify:      *verify,
		AppSecret:   *appSecret,
		VerifyToken: *verifyToken,
		Token:       *pageToken,
	})

	// Setup a handler to be triggered when a message is received
	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		p, err := client.ProfileByID(m.Sender.ID)
		if err != nil {
			fmt.Println("Something went wrong!", err)
		}
		r.Text(fmt.Sprintf("%v: %v", p.FirstName, m.Sender.ID), messenger.ResponseType)
	})

	server := &http.Server{Addr: ":8080", Handler: client.Handler()}
	log.Printf("App is available at :8080")
	server.ListenAndServe()
}
