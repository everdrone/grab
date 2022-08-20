global {
  location = "{{ .Location }}"
}

site "unsplash" {
  test = "unsplash"

  asset "image" {
    pattern = "contentUrl\":\"([^\"]+)\""
    capture = 1

    transform filename {
      pattern = "(?:.+)photos\\/(.*)"
      replace = "$${1}.jpg"
    }
  }

  info "title" {
    pattern = "meta[^>]+property=\"og:title\"[^>]+content=\"(?P<title>[^\"]+)\""
    capture = "title"
  }

  subdirectory {
    pattern = "\\(@(?P<username>\\w+)\\)"
    capture = "username"
    from    = body
  }
}