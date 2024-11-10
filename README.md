# Vigo Bus Telegram Bot

## Prerequisites

- Have the Vigo Bus Core installed on the same namespace (refer to https://github.com/eryalito/vigo-bus-core)

## Installation

```bash
helm upgrade --install tgbot --set botToken=<your-bot-token> oci://ghcr.io/eryalito/vigo-bus-chart/vigo-bus-tg-bot -n <namespace>
```
