package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yeslayla/birdbot-common/common"
)

type feedbackWebhookModule struct {
	webhookURL     string
	payloadType    string
	successMessage string
	failureMessage string
}

type FeedbackWebhookConfiguration struct {
	SuccessMessage string
	FailureMessage string
	PayloadType    string
}

// NewFeedbackWebhookComponent creates a new component
func NewFeedbackWebhookComponent(webhookURL string, config FeedbackWebhookConfiguration) common.Module {
	m := &feedbackWebhookModule{
		webhookURL:     webhookURL,
		payloadType:    "default",
		successMessage: "Feedback recieved!",
		failureMessage: "Failed to send feedback!",
	}

	if config.SuccessMessage != "" {
		m.successMessage = config.SuccessMessage
	}
	if config.FailureMessage != "" {
		m.failureMessage = config.FailureMessage
	}
	if config.PayloadType != "" {
		m.payloadType = config.PayloadType
	}

	return m
}

func (c *feedbackWebhookModule) Initialize(birdbot common.ModuleManager) error {
	birdbot.RegisterCommand("feedback", common.ChatCommandConfiguration{
		Description:       "Sends a feedback message",
		EphemeralResponse: true,
		Options: map[string]common.ChatCommandOption{
			"message": {
				Description: "Content of what you'd like to communicate in your feedback.",
				Type:        common.CommandTypeString,
				Required:    true,
			},
		},
	}, func(user common.User, args map[string]any) string {

		message, ok := args["message"]
		if !ok {
			return "Missing content in command"
		}

		var data []byte

		// Supported payload types
		switch c.payloadType {
		case "discord":
			data, _ = json.Marshal(map[string]any{
				"content": fmt.Sprintf("%s: %s", user.DisplayName, message),
			})
		case "slack":
			data, _ = json.Marshal(map[string]any{
				"text": fmt.Sprintf("%s: %s", user.DisplayName, message),
			})
		default:
			data, _ = json.Marshal(map[string]any{
				"message":  message,
				"username": user.DisplayName,
			})
		}

		body := bytes.NewBuffer(data)

		// Send HTTP request
		resp, err := http.Post(c.webhookURL, "application/json", body)
		if err != nil {
			log.Printf("Failed to post feedback to url '%s': %s", c.webhookURL, err)
			return c.failureMessage
		}

		// Validate response
		if resp.Status[0] != '2' {
			log.Printf("Webhook returned %v: %s", resp.Status, message)
			return c.failureMessage
		}

		// Read body for any special response
		response := map[any]any{}
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.successMessage
		}

		if err = json.Unmarshal(responseBody, &response); err != nil {
			return c.successMessage
		}

		if message, ok := response["message"]; ok {
			v := fmt.Sprint(message)
			if len(v) > 0 {
				return v
			}
		}

		return c.successMessage
	})
	return nil
}
