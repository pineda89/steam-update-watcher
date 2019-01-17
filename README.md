# steam-update-watcher

golang application which sends a telegram message on steam game update

| env var | description | example |
| ------ | ------ | ------ |
| TELEGRAMTOKEN | your telegram bot token | 2321121:AAG6ZZZZZZZZZZZJJJ99h9g |
| TELEGRAMCHATID | telegram chat id to send the expected messages | 12341234 |
| REGEXP | regular expression to match the updates | .*(App)+.*(12578080)+.*(PLAYERUNKNOWN)+.*(BATTLEGROUNDS) |

### About telegram bot

Open telegram, and message to @BotFather to create a new bot

### docker

docker build -t steam-update-watcher .

docker run -d -e TELEGRAMTOKEN="2321121:AAG6ZZZZZZZZZZZJJJ99h9g" -e TELEGRAMCHATID=12341234 -e REGEXP=".*(App)+.*(12578080)+.*(PLAYERUNKNOWN)+.*(BATTLEGROUNDS)" --name steam-update-watcher steam-update-watcher
