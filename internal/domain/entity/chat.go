package entity

import "errors"

type ChatConfig struct {
	Model            *Model
	Temperature      float32
	TopP             float32
	N                int
	Stop             []string
	MaxTokens        int
	PresencePenalty  float32
	FrequencyPenalty float32
}

type Chat struct {
	ID                   string
	UserID               string
	Messages             []*Message
	InitialSystemMessage []*Message
	ErasedMessages       []*Message
	TokensUsage          int
	Status               string
	Config               *ChatConfig
}

func (c *Chat) AddMessage(m *Message) error {
	if c.Status == "ended" {
		return errors.New("chat is ended, no more messages allowed")
	}
	for {
		if c.Config.Model.GetMaxToken() >= m.GetQtdTokens()+c.TokensUsage {
			c.Messages = append(c.Messages, m)
			c.RefreshTokenUsage()
			break
		}
		c.ErasedMessages = append(c.ErasedMessages, c.Messages[0])
		c.Messages = c.Messages[1:]
		c.RefreshTokenUsage()
	}
	return nil
}

func (c *Chat) RefreshTokenUsage() {
	c.TokensUsage = 0
	for m := range c.Messages {
		c.TokensUsage += c.Messages[m].GetQtdTokens()

	}

}
