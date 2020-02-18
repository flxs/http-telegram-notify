# http-telegram-notify

Receives messages via an HTTP endpoint and sends them to Telegram users using the Telegram Bot API.

It's intended to run in Kubernetes and provide other pods with an easy way to send notifications via simple `curl` calls,
while keeping the Telegram API token secret.

## Configuration using environment variables

Required:

- `APP_TOKEN`: A token that needs to be sent in the `X-Auth-Token` header with every request
- `TELEGRAM_TOKEN`: A Telegram Bot API token

Optional:

- `DEBUG`: If set, debug logging is enabled

## Running locally

Starting a `http-telegram-notify` Docker container:

```
docker run --rm -e APP_TOKEN=your-app-token -e TELEGRAM_TOKEN=your-telegram-token -p 8080:8080 http-telegram-notify
```

Message recipients are expected to be given as a Telegram Chat ID; you can find out yours using the [Telegram JSON dump bot](https://t.me/jsondumpbot); 
send it a message and look in the output for the value in `message.from.id`; should be a relatively lengthy numerical value.

To send a message to Telegram Chat ID `1234567890`:

```
curl -XPOST \
    localhost:8080/msg \
    -H "X-Auth-Token: your-app-token" \
    -d '{"to":1234567890,"message":"Hello World!", "silent":false}' 
```

Where:

-  `to` contains the chat ID to send the message to (as an integer)
-  `message` contains the message text (may contain emoji, like `:+1:`)
-  `silent` specifies whether the message should trigger an audible/vibration notification (`true`: no notification, `false`: with notification)

## Running on Kubernetes

This assumes a running Kubernetes cluster and `kubectl` configured locally. The following is intended to help you get something running quickly, but you'll
want to make some improvements before deploying this outside a development cluster.

First, create a namespace, e.g. `telegram`:

```
kubectl create namespace telegram
```

Next, generate a random app token – this will be used by clients to authenticate with `http-telegram-notify`), e.g. using something like `openssl rand -base64 48`.

Then, create secrets for the app token and the Telegram Bot API token:

```
kubectl create secret generic http-telegram-notify-token --namespace telegram --from-literal token="your-random-app-token"
kubectl create secret generic telegram-bot-api-token --namespace telegram --from-literal token="your-telegram-bot-api-token"
```

Create a Service and a Deployment using something like this:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: http-telegram-notify
  namespace: telegram
spec:
  selector:
    app: http-telegram-notify
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: http-telegram-notify
  namespace: telegram
  labels:
    app: http-telegram-notify
spec:
  replicas: 1
  selector:
    matchLabels:
      app: http-telegram-notify
  template:
    metadata:
      labels:
        app: http-telegram-notify
    spec:
      containers:
      - name: http-telegram-notify
        image: 1flx/http-telegram-notify:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_TOKEN
          valueFrom:
            secretKeyRef:
              name: http-telegram-notify-token
              key: token
        - name: TELEGRAM_TOKEN
          valueFrom:
            secretKeyRef:
              name: telegram-bot-api-token
              key: token
```

To test the service, forward port 80 on the service to 8080 on your local machine:

```
kubectl port-forward --namespace telegram service/http-telegram-notify 8080:80
```

Try to send a message:

```
curl -XPOST \
    localhost:8080/msg \
    -H "X-Auth-Token: your-app-token" \
    -d '{"to":1234567890,"message":"Hello World!", "silent":false}' 
```

## Building locally

To build an `alpine`-based Docker image with `http-telegram-notify` installed, run:

```
docker build . -t http-telegram-notify
```