package consts

const (
	StartMsg = "Hi! 👋🥸 This is a Telegram bot for communicating with chatGpt language model." +
		"\n\n⚡️ You have access to 5 queries absolutely free of charge daily. (They are refreshed every 24 hours)." +
		"\n\n💎💎💎 You can also invite your friend and get a free \"Basic\" plan for a week. To get a referral link, use the command /invite" +
		"\n\nIf you need more requests, we have several options for monthly subscriptions:\n\n🤓 Basic Plan: 25 requests per day - 5$\n⚜️ VIP Plan: 50 requests per day - 10$\n🇦🇪 DELUXE Plan: 100 requests per day - 15$" +
		"\n\n💸 Buy subscription: /buy"
	BuyMsg = "Monthly subscriptions:\n\n🤓 Basic Plan: 25 requests per day - 5$\n⚜️ VIP Plan: 50 requests per day - 10$\n🇦🇪 DELUXE Plan: 100 requests per day - 15$" +
		"\n\n💸 Buy subscription: \nTether(USDT) TRC20: TBbdyory2Bc4csT5Xc4tieEHptD2XEKeRW" +
		"\n\n💬 Contact: \n@dendefoe"
	ErrorMsg  = "❌ Something wen't wrong. 🤕 Try again later."
	StatusMsg = "⚡️ Your subscription:\n%s (%d requests per day) \n💫 Available: %d requests\n🕐 Valid due: %s\n\n💬 Contact: \n@dendefoe" +
		"\n\n💚 Invite your friend and get a free \"Basic\" plan for a week. To get a referral link, use the command /invite" +
		"\n\n💸 Buy subscription: /buy"
	MissingGptKey            = "❌ missing variable: GPT_API_KEY"
	MissingTgKey             = "❌ missing variable: TELEGRAM_API_KEY"
	UnsupportedMessageType   = "❌ sorry, the message type you sent is not supported yet"
	MessageIsTooShort        = "❌ you have sent a message that is too short. The minimum number of characters is 4"
	SubscriptionPlanHomeless = "🗿 Homeless"
	SubscriptionPlanBasic    = "🤓 Basic"
	SubscriptionPlanVIP      = "⚜️ VIP"
	SubscriptionPlanDeluxe   = "🇦🇪 DELUXE"
	SubscriptionPlanHacker   = "🦄 Hacker"
	ModeSuccess              = "💚 Mode %s"
	DialogSuccess            = "🧠 Dialog mode: %s"
	RunOutOfMessages         = "🦄 Sorry, you ran out of messages\n\n💬 Contact: \n@dendefoe\n\n💚 Invite your friend and get a free \"Basic\" plan for a week. To get a referral link, use the command /invite\n\n💸 Buy subscription: /buy"
	InviteLink               = "💚 Share with your friend and get the 🤓 Basic Plan for a week: \n🔗 https://t.me/call_fresco_bot?start=ref%d"
	SuccessRef               = "Congratulations! 🎉 Your referral was successful, and a new user has joined through your link! As a token of our appreciation, we have activated your bonus subscription. 🎁 Enjoy the extended features and thank you for spreading the word! 👍 If you have any questions, feel free to ask. Happy chatting! 🤖💬"
)
