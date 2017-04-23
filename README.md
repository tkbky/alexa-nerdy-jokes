# Alexa Nerdy Jokes

Alexa that tells nerdy jokes.

## Prerequisites

1. Sign up for an AWS developer account. [link](https://developer.amazon.com/)
2. Create an Alexa skill. [link](https://developer.amazon.com/edw/home.html#/skill/create/)

## Integration Testing

Firstly, get a public url using [ngrok](https://ngrok.com/) or [localtunnel](https://localtunnel.github.io/www/) or anything that can expose your local to the world.

Then, on AWS Alexa skill console,

1. Under **Configuration**, set the **Service Endpoint to your local entry point**, e.g. `https://xxxxxxxx.ngrok.io/echo/nerdyjokes`
2. Under **Test**, there is a **Service Simulator**. Enter an utterance, then click **Ask Nerdy Jokes**.

## Heroku Deployment

1. Install required add-ons

  ```bash
  heroku addons:create heroku-postgresql:hobby-dev
  ```

2. Set the environment variables

  ```bash
  heroku config:set ALEXA_SKILL_APP_ID=amzn1.ask.skill.xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  heroku config:set GOVERSION=go1.8 # Optional, heroku will set the default to latest it has
  heroku config:set GO_INSTALL_PACKAGE_SPEC=./ # Optional
  heroku config:set PORT=5000
  ```
