global {
  location = "{{.Location}}"

  network {
    timeout = 30
  }
}

site "unsplash" {
  test = "unsplash"

  network {
    inherit = false
    retries = 3
    headers = {
      "User-Agent" : "grab/0.1",
    }
  }

  asset "image" {
    pattern = "contentUrl\":\"([^\"]+)\""
    capture = 1

    network {
      inherit = false
    }

    transform url {
      pattern = "([^/]+)$"
      replace = "https://unsplash.com/$1"
    }

    transform filename {
      pattern = "(https?[^\\?\\#]+)(.*)"
      replace = "$${1}.jpg"
    }
  }

  info "title" {
    pattern = "meta[^>]+property=\"og:title\"[^>]+content=\"(?P<title>[^\"]+)\""
    capture = "title"
  }

  subdirectory {
    pattern = "href=\"\\/\\@(?P<user>[^\"]+)"
    capture = "user"
    from    = body
  }
}