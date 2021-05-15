---
title: "Local Dev Setup with automatic HTTPS + WebSockets in 1 Line"
date: 2021-02-02T16:11:13Z
draft: false
canonicalUrl: "https://dev.to/remast/using-sonarcloud-with-github-actions-and-maven-31kg"
---

Let's set up a local web server for development that supports automatic HTTPs and WebSockets in just 1 line of code. No problem thanks to the awesome [Caddyserver](https://caddyserver.com/).

## Preparations

### Set up the local hostname
For automatic HTTPs you need to set up a local hostname first. We will use `remast.local`. To make your local development machine available at that hostname you need to add the hostname to your local hosts file.

So locate the hosts file on your machine and add the following line:
```
127.0.0.1	remast.local
```

Here's where you find the hosts file in your operating system.

| Operating System | Location of hosts file           |
| ---------------- |:-------------|
| Windows 10       | `C:\Windows\System32\drivers\etc\hosts` |
| Linux       | `/etc/hosts` |
| Mac OS X       | `/private/etc/hosts` |

Now verify by testing to reach the hostname via ping with `ping remast.local`. If that takes you to `127.0.0.1` everything is fine and you can proceed.


### Install Caddyserver

Now install [Caddyserver](https://caddyserver.com/) on your machine. If you're on Windows and use the [Chocolatey package manager](https://chocolatey.org) you can do that with `choco install caddy`. For other plattforms check the [Caddyserver install docs](https://caddyserver.com/docs/install).

## Starting the local dev server with HTTPs and WebSockets

Now off we go with our local dev server with automatic HTTPs and WebSocket support. Here's the only one line you need:

```
caddy reverse-proxy --from remast.local --to 127.0.0.1:8080
```

This one line will fire up your local [Caddyserver](https://caddyserver.com/) which will automaticlally start HTTPs and proxy all requests including WebSockets. How cool is that?!

### Using the Caddyfile

Of course you can always save your configuration in file called `Caddyfile`, with the following contents:
```
remast.local {
	reverse_proxy 127.0.0.1:8080
}
```

### What's next?

Once you've made your first steps [Caddyserver](https://caddyserver.com/) there's no way back. Caddy is really awesome for both development and production!