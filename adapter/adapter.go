package adapter

import (
	"chat/adapter/baichuan"
	"chat/adapter/bing"
	"chat/adapter/claude"
	"chat/adapter/dashscope"
	"chat/adapter/hunyuan"
	"chat/adapter/midjourney"
	"chat/adapter/oneapi"
	"chat/adapter/palm2"
	"chat/adapter/skylark"
	"chat/adapter/slack"
	"chat/adapter/zhinao"
	"chat/adapter/zhipuai"
	"chat/globals"
	"chat/utils"
)

type ChatProps struct {
	Model    string
	Plan     bool
	Infinity bool
	Message  []globals.Message
	Token    int
}

func NewChatRequest(props *ChatProps, hook globals.Hook) error {
	if oneapi.IsHit(props.Model) {
		return createRetryOneAPI(props, hook)

	} else if globals.IsChatGPTModel(props.Model) {
		return createRetryChatGPTPool(props, hook)

	} else if globals.IsClaudeModel(props.Model) {
		return claude.NewChatInstanceFromConfig().CreateStreamChatRequest(&claude.ChatProps{
			Model:   props.Model,
			Message: props.Message,
			Token:   utils.Multi(props.Token == 0, 50000, props.Token),
		}, hook)

	} else if globals.IsSparkDeskModel(props.Model) {
		return retrySparkDesk(props, hook, 0)

	} else if globals.IsPalm2Model(props.Model) {
		return palm2.NewChatInstanceFromConfig().CreateStreamChatRequest(&palm2.ChatProps{
			Model:   props.Model,
			Message: props.Message,
		}, hook)
	} else if globals.IsSlackModel(props.Model) {
		return slack.NewChatInstanceFromConfig().CreateStreamChatRequest(&slack.ChatProps{
			Message: props.Message,
		}, hook)
	} else if globals.IsBingModel(props.Model) {
		return bing.NewChatInstanceFromConfig().CreateStreamChatRequest(&bing.ChatProps{
			Model:   props.Model,
			Message: props.Message,
		}, hook)
	} else if globals.IsZhiPuModel(props.Model) {
		return zhipuai.NewChatInstanceFromConfig().CreateStreamChatRequest(&zhipuai.ChatProps{
			Model:   props.Model,
			Message: props.Message,
		}, hook)
	} else if globals.IsQwenModel(props.Model) {
		return dashscope.NewChatInstanceFromConfig().CreateStreamChatRequest(&dashscope.ChatProps{
			Model:   props.Model,
			Message: props.Message,
		}, hook)
	} else if globals.IsMidjourneyModel(props.Model) {
		return midjourney.NewChatInstanceFromConfig().CreateStreamChatRequest(&midjourney.ChatProps{
			Model:    props.Model,
			Messages: props.Message,
		}, hook)
	} else if globals.IsHunyuanModel(props.Model) {
		return hunyuan.NewChatInstanceFromConfig().CreateStreamChatRequest(&hunyuan.ChatProps{
			Model:    props.Model,
			Messages: props.Message,
		}, hook)
	} else if globals.Is360Model(props.Model) {
		return zhinao.NewChatInstanceFromConfig().CreateStreamChatRequest(&zhinao.ChatProps{
			Model:   props.Model,
			Message: props.Message,
			Token:   utils.Multi(props.Token == 0, 2048, props.Token),
		}, hook)
	} else if globals.IsBaichuanModel(props.Model) {
		return baichuan.NewChatInstanceFromConfig().CreateStreamChatRequest(&baichuan.ChatProps{
			Model:   props.Model,
			Message: props.Message,
			Token:   utils.Multi(props.Token == 0, 4096, props.Token),
		}, hook)
	} else if globals.IsSkylarkModel(props.Model) {
		return skylark.NewChatInstanceFromConfig().CreateStreamChatRequest(&skylark.ChatProps{
			Model:   props.Model,
			Message: props.Message,
			Token:   utils.Multi(props.Token == 0, 4096, props.Token),
		}, hook)
	}

	return hook("Sorry, we cannot find the model you are looking for. Please try another model.")
}
