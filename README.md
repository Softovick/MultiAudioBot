# MultiAudioBot

STEPS FOR RUN LOCAL ENVIRONMENT

1. Set environment variable TOKEN_TELEGRAM. Value get from [Telegram bot documentation](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
   * Windows Powershell  <code>env:TOKEN_TELEGRAM = "sometoken";</code>
   * Windows cmd <code>TOKEN_TELEGRAM=sometoken</code>
   * Linux bash <code>export TOKEN_TELEGRAM=sometoke</code>
2. Build bot.go and run resulting binary file (bot/bot.exe)

STEPS FOR RUN DOCKER ENVIRONMENT

1. Make .env, write TOKEN_TELEGRAM=sometoken
2. docker-compose up