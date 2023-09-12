# Glean

A plain and simple blog engine in Go. It serves Markdown files from a folder
as HTML with a little bit of styling to make it all bearable and nice.

## Usage

First, set the following environment variables:

```
GLEAN_PORT="8080"
GLEAN_PATH="path/to/folder"
GLEAN_PASS="biglongsecureapikey"
GLEAN_TITLE="Joe Schmoe's Blog"
GLEAN_AUTHOR="Joe Schmoe"
GLEAN_EMAIL="joe@schmoe.com"
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
- Configurable path for folder containing post files
- RSS feed
