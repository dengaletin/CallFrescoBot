package consts

const (
	StartMsg = "Hi! 👋🥸 This is a Telegram bot for communicating with chatGpt language model." +
		"\n\n⚡️ You have access to 15 queries absolutely free of charge daily. (They are refreshed every 24 hours)." +
		"\n\n If you need more requests, we have several options for monthly subscriptions:\n\n🤓 Student: 50 requests per day - 99 RUB\n⚜️ Major: 100 requests per day - 250 RUB\n🇦🇪 GigaSheikh: 200 requests per day - 399 RUB"
	ErrorMsg                   = "❌ Something wen't wrong. 🤕 Try again later."
	StatusMsg                  = "⚡️ Your subscription:\n%s (%d requests per day) \n💫 Available: %d requests\n\n💬 Contact: \n@dendefoe\n💸 Buy subscription: \nhttps://www.donationalerts.com/r/dendefoe"
	MissingGptKey              = "❌ Missing variable: GPT_API_KEY"
	MissingTgKey               = "❌ Missing variable: TELEGRAM_API_KEY"
	UnsupportedMessageType     = "❌ Sorry, the message type you sent is not supported yet."
	MessageIsTooShort          = "❌ You have sent a message that is too short. The minimum number of characters is 4."
	SubscriptionPlanBomj       = "🗿 Bomj"
	SubscriptionPlanStudent    = "🤓 Student"
	SubscriptionPlanMajor      = "⚜️ Major"
	SubscriptionPlanGigaSheikh = "🇦🇪 GigaSheikh"
	SubscriptionPlanHacker     = "🦄 Hacker"
	RunOutOfMessages           = "🦄 Sorry, you ran out of messages\n\n💬 Contact: \n@dendefoe\n💸 Buy subscription: \nhttps://www.donationalerts.com/r/dendefoe"
)
