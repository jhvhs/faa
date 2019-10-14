# faa

The simplest Slack retrobot you ever did see

## Deploying

### Slack Integration

FAA must be configured in your Slack as a "Custom Integration". Here's an example of the production integration settings.

![faa_slack_integration](assets/faa_slack_integration.png)

Make sure to configure `URL` to the publicly available URL of your Cloud Foundry app (see below)


### Cloud Foundry

FAA runs as an app on Cloud Foundry. To successfully push, you must provide the following:

- `SLACK_VERIFICATION_TOKEN`: *string*, Verification token provided by your slack integration, see "Token" in the slack configuration above
- `POSTFACTO_RETRO_ID`: *integer*, The postfacto ID of your regular retro

It's also possible to provide the following optional settings:
- `POSTFACTO_APP_URL`: *string*, The URL of the main browser app, defaults to `https://retros.cfapps.io`
- `POSTFACTO_API_URL`: *string*, The URL of the postfacto API, defaults to `https://retros-iad-api.cfapps.io`
- `POSTFACTO_RETRO_PASSWORD`: *string*, Set this when the retro board is password-protected

Other configuration necessary to run on Cloud Foundry can be found in our [production manifest.yml](manifest.yml)


## Using

Assuming you have configured your slack integration with the command `/retro`

```
/retro [happy/meh/sad/tech] [your message]
```


## Development

- Uses [gvt](github.com/FiloSottile/gvt) for vendoring
- Convenient `./bin/build` script