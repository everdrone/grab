package testutils

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
)

const htmlPage string = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Grab Test Server</title>
</head>
<body>
  <h1>Grab Test Server</h1>
  <p>Author: @everdrone</p>
  <p>url: <a href="https://github.com/everdrone/grab">GitHub Repo</a></p>
  <div>
    <img src="{{ .Base }}/img/a.jpg" />
    <img src="{{ .Base }}/img/b.jpg" />
    <img src="{{ .Base }}/img/c.jpg" />
  </div>
  <div>
    <img src="/img/a.jpg" />
    <img src="/img/b.jpg" />
    <img src="/img/c.jpg" />
  </div>
  <div>
    <video src="{{ .Base }}/video/a/small.mp4" ></video>
    <video src="{{ .Base }}/video/b/small.mp4" ></video>
    <video src="{{ .Base }}/video/c/small.mp4" ></video>
  </div>
  <div>
    <audio controls>
      <source src="{{ .Base }}/audio/a/preview" type="audio/mpeg">
      <source src="{{ .Base }}/audio/b/preview" type="audio/ogg">
      <source src="{{ .Base }}/audio/c/preview" type="audio/aac">
    </audio>
  </div>
  <div>
    <img src="{{ .Base }}/secure/a.jpg" />
    <img src="{{ .Base }}/secure/b.jpg" />
    <img src="{{ .Base }}/secure/c.jpg" />
  </div>
  <div>
    <img src="{{ .Base }}/not-found/a.jpg" />
    <img src="{{ .Base }}/not-found/b.jpg" />
    <img src="{{ .Base }}/not-found/c.jpg" />
  </div>
  <a href="{{ .Base }}/broken/1.jpg">absolutely broken</a>
  <a href="/broken/2.jpg">relatively broken</a>
</body>
</html>`

func CreateMockServer() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.GET("/gallery/:id/:type", func(c echo.Context) error {
		if c.Param("type") == "test" && c.Param("id") == "123" {
			buf := new(bytes.Buffer)

			addr := "http://" + strings.Replace(e.ListenerAddr().String(), "[::]", "127.0.0.1", -1)

			page := template.Must(template.New("test").Parse(htmlPage))
			if err := page.Execute(buf, map[string]string{"Base": addr}); err != nil {
				log.Fatalf("error executing template: %v", err)
			}

			return c.HTML(http.StatusOK, buf.String())
		}
		return c.NoContent(http.StatusNotFound)
	})

	e.GET("/broken/:id", func(c echo.Context) error {
		// will cause a reading error
		c.Response().Header().Set("Content-Length", "999")
		return c.NoContent(http.StatusOK)
	})

	e.GET("/img/:id", func(c echo.Context) error {
		for _, id := range []string{"a", "b", "c"} {
			if c.Param("id") == id+".jpg" {
				return c.String(http.StatusOK, "image"+id)
			}
		}

		return c.NoContent(http.StatusNotFound)
	})

	e.GET("/video/:id/:size", func(c echo.Context) error {
		for _, size := range []string{"small", "large"} {
			for _, id := range []string{"a", "b", "c"} {
				if c.Param("id") == id && c.Param("size") == size+".mp4" {
					return c.String(http.StatusOK, "video"+id+size)
				}
			}
		}

		return c.NoContent(http.StatusNotFound)
	})

	e.GET("/audio/:id/:fileType", func(c echo.Context) error {
		for _, fileType := range []string{"small", "large"} {
			for _, id := range []string{"a", "b", "c"} {
				if c.Param("id") == id && c.Param("fileType") == fileType {
					return c.String(http.StatusOK, "audio"+id+fileType)
				}
			}
		}

		return c.NoContent(http.StatusNotFound)
	})

	e.GET("/secure/:id", func(c echo.Context) error {
		if c.Request().Header.Get("custom_header") == "123" {
			for _, id := range []string{"a", "b", "c"} {
				if c.Param("id") == id+".jpg" {
					return c.String(http.StatusOK, "secure"+id)
				}
			}
		}

		return c.NoContent(http.StatusUnauthorized)
	})

	// used to get the server address during execution.
	// on windows, this triggers the firewall to open the port.
	l, _ := net.Listen("tcp4", "127.0.0.1:0")
	e.Listener = l

	return e
}

/*
# configuration for this server

global {
	location =
}

site "example" {
	test = "http:\\/\\/127\\.0\\.0\\.1:\\d+"

	asset "image" {
		pattern = "<img src=\"([^\"]+/img/[^\"]+)"
		capture = 1
		find_all = true
	}

	asset "video" {
		pattern = "<video src=\"([^\"]+)"
		capture = 1
		find_all = true

		transform url {
			pattern = "(.+)small(.*)"
			replace = "$${1}large$2"
		}

		transform filename {
			pattern = "\\/video\\/(?P<id>\\w+)\\/(\\w+)\\.(?P<extension>\\w+)"
			replace = "$${id}.$${extension}"
		}
	}

	asset "audio" {
		pattern = "<source src=\"(?P<named_group>[^\"]+)"
		capture = "named_group"
		find_all = true
	}

	asset "secure" {
		pattern = "<img src=\"([^\"]+/secure/[^\"]+)"
		capture = 1
		find_all = true

		network {
			headers = {
				"custom_header" = "123"
			}
		}
	}

	info "author" {
		pattern = "Author: @(?P<username>[^<]+)"
		capture = "username"
	}

	info "title" {
		pattern = "<title>([^<]+)"
		capture = 1
	}

	subdirectory {
		pattern = "\\/gallery\\/(\\d+)"
		capture = 1
		from = url
	}
}
*/
