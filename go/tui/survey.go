package tui

import (
	"github.com/AlecAivazis/survey/v2"
)

func Survey(message string, options []string) string {
  var input string
  prompt := &survey.Select{
    Message: message,
    Options: options,
  }
  survey.AskOne(prompt, &input)
  return input
}
