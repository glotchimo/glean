# Glean

A plain and simple blog engine in Go. It serves Markdown files from a folder
as HTML with a little bit of styling to make it all bearable and nice.

## Setup

To get started, you'll need a config like this:

    host: http://localhost:8080
    posts_path: ../posts
    emails_path: emails.txt
    
    meta:
      title: My Awesome Blog
      author: Me, the Author
      email: author@blog.url
      links:
        Git: https://github.com/me-the-author
        YouTube: https://www.youtube.com/channel/me-the-author
    
    smtp:
      host: smtp.mail.com
      port: 587
      username: author@blog.url
      password: p4ssw0rd
      sender: author@blog.url

You don't have to set up an SMTP server if you don't want; it'll still work
without it, you'll just see some errors in the server logs.

## Usage

To launch your blog, use the CLI:

    glean -conf prod.yml -port 80

As long as your setup is good, it'll get going and you're in business!

## Features

- [x] On-the-fly Markdown -> HTML conversion for posts
- [x] Dynamically-generated index page
- [x] Configurable header content
- [x] Configurable path for folder containing post files
- [x] Mailer & form for new post notifications
- [ ] SQLite database for storing mailing list
- [ ] Dynamically-filled SLOC count in footer
