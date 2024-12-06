# Glean

A plain and simple blog engine in Go. It serves Markdown files from Tigris object storage
with a little bit of styling to make it all bearable and nice.

## Usage

First, set the following environment variables:

```
GLEAN_PORT="8080"
GLEAN_PASS="biglongsecureapikey"
GLEAN_TITLE="Joe Schmoe's Blog"
GLEAN_AUTHOR="Joe Schmoe"
GLEAN_EMAIL="joe@schmoe.com"

# Tigris Storage Configuration
BUCKET_NAME="summer-grass-2004"
AWS_ENDPOINT_URL_S3="https://fly.storage.tigris.dev"
AWS_ACCESS_KEY_ID="tid_xxxxxx"
AWS_SECRET_ACCESS_KEY="tsec_xxxxxx"
```

Then, you're good to run it:

```sh
./glean
```

Use the `post` script to add posts:

```sh
post -d Article.md
```

## Features

- On-the-fly Markdown -> HTML conversion for posts
- Dynamically-generated index page
- Configurable header content
- Tigris object storage for posts
- RSS feed
