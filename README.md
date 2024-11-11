# Vigo Bus Telegram Bot

Vigo Bus Telegram Bot (unofficial). Get routes, stops and real time info about the buses in the city of Vigo. The bot is currently running on https://t.me/busvigobot or @busvigobot. 

## Features

- Favorite stops
- Routes
- Stops
- Real time schedules

## Prerequisites

- Have the Vigo Bus Core installed on the same namespace (refer to https://github.com/eryalito/vigo-bus-core)

## Installation

```bash
helm upgrade --install tgbot --set botToken=<your-bot-token> oci://ghcr.io/eryalito/vigo-bus-chart/vigo-bus-tg-bot -n <namespace>
```
