package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetResponse(client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))

	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		panic("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				for i := 0; i < 50; i++ {
					time.Sleep(2 * time.Millisecond)
					fmt.Print("-")
				}
				time.Sleep(15 * time.Millisecond)
				fmt.Println()
				fmt.Print(" 🚀Ask me anything ('quit' to end): ")
				if !scanner.Scan() {
					break
				}

				question := scanner.Text()
				switch question {
				case "quit":
					fmt.Println("  Shutting down AI-CLI")
					time.Sleep(1 * time.Second)
					for i := 0; i < 5; i++ {
						time.Sleep(500 * time.Millisecond)
						fmt.Print("<><><><><>")
					}
					fmt.Println()
					fmt.Println("  See you next time!")
					time.Sleep(250 * time.Millisecond)
					fmt.Println("Developed by Anish Dubey ")
					time.Sleep(250 * time.Millisecond)
					fmt.Println("https://twitter.com/anish_dubey_")
					time.Sleep(250 * time.Millisecond)
					for i := 0; i < 50; i++ {
						time.Sleep(2 * time.Millisecond)
						fmt.Print("-")
					}
					quit = true

				default:
					GetResponse(client, ctx, question)
				}
			}
		},
	}

	rootCmd.Execute()
}
