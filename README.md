# tagesschau

A small command line utility to conveniently watch "Tagesschau".

This tool searches for the right video, applies some custom playback
settings<sup>1</sup> and then plays it.

<sup>1</sup>*Playback Settings*
* Skip Intro
* 2x video speed

## Usage

```
Usage: tagesschau [options] [url]
  -d string
      Target date of the video (default <current date>)
  -o int
      Date offset in days
  <url>
      If an url is given, the linked video will be played. Any date options will be ignored.
```


## Dependencies

Because this tool is a wrapper, it has the following external dependencies:

* `youtube-dl`
* `mpv`
