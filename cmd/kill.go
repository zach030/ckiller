package cmd

import (
	"context"
	"encoding/json"

	"fmt"
	"log"
	"os/exec"

	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// askCmd ask ck for help
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "ask ck for help",
	Run:   askCmdFunc,
}

const PromptTmp = `I want you to act as an experienced developer，
I will provide you with the frequently-used command and my question about that, you can refer the given command help document.
My question is wrapped by {}, If you found the question and command are not related, just answer "I don't know". 
You can think step by step, and provide the answer in JSON Array format with the following keys: step_id, action, command.
action means the comment of the command you provide, command refers to the plain command user can run directly without any comment. 
If the command is empty, just keep it blank.`

func init() {
	RootCmd.AddCommand(askCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Step struct {
	StepID  int64  `json:"step_id"`
	Action  string `json:"action"`
	Command string `json:"command"`
}

func askCmdFunc(cmd *cobra.Command, args []string) {
	help := commandHelp(args[0])
	if cfg.OpenaiApiKey == "" {
		log.Panic("empty apikey")
	}
	openaiCfg := openai.DefaultConfig(cfg.OpenaiApiKey)
	if cfg.OpenaiBaseURL != "" {
		openaiCfg.BaseURL = cfg.OpenaiBaseURL
	}
	var prompt = PromptTmp
	prompt = preparePrompt(cfg.Language)
	client := openai.NewClientWithConfig(openaiCfg)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("command:%s, my question is:{%s}, help doc:%s,", args[0], args[1], help),
				},
			},
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	var content string
	for _, choice := range resp.Choices {
		content += choice.Message.Content
	}
	steps := make([]Step, 0)
	if err = json.Unmarshal([]byte(content), &steps); err != nil {
		log.Panic(err)
		return
	}
	for _, step := range steps {
		color.Green("🚀Step-%d\n", step.StepID)
		color.Cyan("\tAction: %s\n\tCommand: %s\n", step.Action, step.Command)
	}
}

func preparePrompt(language string) string {
	if language == "" {
		language = "English"
	}
	return PromptTmp + "\nNotice! You should reply to me in language: " + language
}

func commandHelp(target string) string {
	// create a new *Cmd instance
	// here we pass the command as the first argument and the arguments to pass to the command as the
	// remaining arguments in the function
	cmd := exec.Command(target, "-h")

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, err := cmd.Output()
	if err != nil {
		// if there was any error, print it here
		fmt.Println("could not run command: ", err)
	}
	return string(out)
}

func manDoc(target string) string {
	// create a new *Cmd instance
	// here we pass the command as the first argument and the arguments to pass to the command as the
	// remaining arguments in the function
	cmd := exec.Command("man", target)

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, err := cmd.Output()
	if err != nil {
		// if there was any error, print it here
		fmt.Println("could not run command: ", err)
	}
	return string(out)
}
