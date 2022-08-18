# Scraping with GRAB

To understand how to exploit the power of Grab, let's use an example.

Say we have a website `https://example.com` that hosts galleries of pictures and videos curated by users. We want to download all the images and video files in the gallery, and we also want to save some text information, like the username of the curator.

An example gallery, located at `https://example.com/gallery/1337` would look something like this:

```html
<html>
  <head>
    <title>My awesome gallery</title>
  </head>
  <body>
    <h1>My awesome gallery</h1>
    <p>
      Curated by <a href="https://example.com/user/everdrone">@everdrone</a>
    </p>
    <p>Wednesday, August 17, 2022</p>
    <div>
      <h1>images</h1>
      <img src="https://cdn.example.com/img/jpg/94257478745" />
      <img src="https://cdn.example.com/img/jpg/20239846093" />
      <img src="https://cdn.example.com/img/jpg/39808447626" />
      <video
        src="https://cdn.example.com/video/1337/39808447626_small.mp4"
      ></video>
    </div>
  </body>
</html>
```

To get the gallery curator we could parse the url to get whatever comes after `user/`. We also have the title both inside `<title>` and inside `<h1>`, the video and image urls are also easily accessible.

## The basics

Let's start by generating a new configuration file using `grab config generate`. Remove everything inside the new `grab.hcl` file and write the following:

```hcl
global {
  location = "/home/<username>/Downloads/grab"
}
```

Replace the path with any directory you want, we will refer to this path as `global.location`.

This will tell grab where to store everything: sites, images, videos, information.

We will add a `site` block called `example`. We will add an attribute `test` with the regular expression that will check if we actually are on `https://example.com`.  
In there we also want to have two `asset` blocks, one for images and one for videos.

```
global {
  location = "/home/<username>/Downloads/grab"
}

site "example" {
  test = ":\\/\\/example\\.com"

  asset "image" {
    pattern  = "<img\\ssrc=\"([^\"]+)"
    capture  = 1
    find_all = true
  }

  asset "video" {
    pattern  = "<video\\ssrc=\"(?P<video_url>[^\"]+)"
    capture  = "video_url"
    find_all = true
  }
}
```

If you're familiar with Regular Expressions, the patterns above should be pretty easy to understand, but we'll go through them here anyway.  
The `asset[image].pattern` expression captures whatever comes after `<img src="` until it finds another double quote, so it will get us the entire image url.  
The expression in `asset[video].pattern` does the same thing but with a `video` tag, and it uses a named capture group, for the sake of this example.

Notice how the `capture` attribute is a number when the regular expression doesn't use named captures, and it's a string in the `video` asset because we want to capture the `(?P<video_url>)` group.

We also want set `find_all` to `true` because we expect to find multiple video and image urls on the site page.

Now let's run the program and see what we get.

```
grab get https://example.com/gallery/1337
```

Internally Grab stores the following urls extracted with the regular expressions defined in each `asset` block:

```
https://cdn.example.com/img/jpg/94257478745
https://cdn.example.com/img/jpg/20239846093
https://cdn.example.com/img/jpg/39808447626
https://cdn.example.com/video/1337/39808447626_small.mp4
```

After Grab finishes downloading everything, we can browse the `global.location` directory to find out that a new `example` directory has been created. This directory is named after the `site "example"` block.

If we navigate inside the `example` directory we will see the following files:

```
94257478745
20239846093
39808447626
39808447626_small.mp4
```

The files without extension are the images, and then we have one video file.

## Subdirectories

Let's organize our downloads by making Grab create subdirectories, so that for any other gallery than the one located at `https://example.com/gallery/1337`, we get a directory named with the gallery id.

```hcl
global {
  location = "/home/<username>/Downloads/grab"
}

site "example" {
  test = ":\\/\\/example\\.com"

  asset "image" {
    pattern  = "<img\\ssrc=\"([^\"]+)"
    capture  = 1
    find_all = true
  }

  asset "video" {
    pattern  = "<video\\ssrc=\"(?P<video_url>[^\"]+)"
    capture  = "video_url"
    find_all = true
  }

  subdirectory {
    pattern = "gallery\\/(\\d+)"
    capture = 1
    from    = url
  }
}
```

By adding a `subdirectory` block inside our site block, we tell grab to create a new directory with the name extracted from the url, by capturing the first group of the expression `gallery\\/(\\d+)`, so in this case, all the assets will be located at `<global.location>/example/1337`.

## Substitutions

There are still a few issues to uncover:

- The images have no extension, but the url suggests that they are encoded using the `jpeg` format.
- The video url ends with `_small`, but we know that there's a better quality option ending with `_large`.

So we will update our configuration file to transform the filename of the `image` assets by moving the `/jpg/` that comes before the name into the extension name. And for the `video` block, we need to transform the url but replacing `_small` with `_large`, before the file is downloaded.

To achieve this, Grab offers the `transform` block.

```hcl
global {
  location = "/home/<username>/Downloads/grab"
}

site "example" {
  test = ":\\/\\/example\\.com"

  asset "image" {
    pattern  = "<img\\ssrc=\"([^\"]+)"
    capture  = 1
    find_all = true

    transform filename {
      pattern = "(https?[^\\?\\#]+)(.*)"
      replace = "$${1}.jpg"
    }
  }

  asset "video" {
    pattern  = "<video\\ssrc=\"(?P<video_url>[^\"]+)"
    capture  = "video_url"
    find_all = true

    transform url {
      pattern =
      replace =
    }
  }

  subdirectory {
    pattern = "gallery\\/(\\d+)"
    capture = 1
    from    = url
  }
}
```
