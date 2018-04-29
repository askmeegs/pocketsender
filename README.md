# pocketsender

Sends your unread Pocket articles to Kindle.


### How to run

1. Get a [Pocket API consumer key](https://getpocket.com/developer/docs/authentication) and access token.

2. Create a local file, `config.json` with these fields populated:

```
{
  "PocketUsername": "",
  "PocketConsumerKey": "",
  "PocketAccessToken": "",
  "FromEmail": "",
  "FromEmailPassword": "",
  "KindleEmail": ""
}
```

3. `docker run -v <path-to-config.json>:/pocketsender/config.json  meganokeefe/pocketsender:0.0.1 --config /pocketsender/config.json`

### Notes 

- Note that right now, pocketsender requires a Google password passed as config -- this is insecure and [will be](https://github.com/m-okeefe/pocketsender/issues/3) addressed.
- pocketsender currently only works with Gmail addresses --> Kindle addresses
- pocketsender sends every unread pocket saved article as a separate Kindle document.
- Document conversion may be lossy depending on the article (as pocketsender converts HTML to PDF, and Kindle converts PDF to kindle format.
- once pocketsender emails a Pocket article, it then archives that article in Pocket, so it will no longer appear in your Pocket "My List."
- when pocketsender runs, it checks for all unread articles, sends them to kindle, then shuts down. (Does not watch your pocket queue or stay running)
- Due to an OpenSSL bug for MacOS, pocketsender currently only runs in linux/Docker.
- If you choose to run in linux from source, note that [wkhtmltopdf](https://wkhtmltopdf.org/) must be pre-installed.
