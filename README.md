# CORS reverse proxy

This is a simple, CORS-enabled reverse proxy. Based off either the flags provided
or the environment variables set, you can use it to proxy a request to an API endpoint
and add CORS headers, so you can make `XMLHTTPRequests`.

To start, launch `proxy-cors` with either the flags `-port` and `-url` or the
environment variables `CORS_PORT` and `CORS_URL`. The default port is `8889` and
the default proxied URL is `http://localhost/`.