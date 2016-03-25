# TFL Status Announcements
## Work in progress

Probably not a thing you should be using. Hack project because I wanted to try out a few things, and like transport data ¯\\\_(ツ)\_/¯

Basic idea: Pull down the TFL service status updates, parse them and convert them to audio announcements with Ivona cloud.

To use;

First configure the relevent APIs, get keys for the [TFL Public Data Feeds](https://api-portal.tfl.gov.uk) and [IVONA Cloud](https://ivona.com/), and add them to the `config.yml` file copied from the example.

You'll also need [Glide](https://github.com/Masterminds/glide) and [Go-Bindata](https://github.com/jteeuwen/go-bindata) installed and in your `$PATH`

To build;

```
$ glide install
$ make
```

Then to run;

```
./tfl-announcements
```

It will log out the status generated.

You can start a web interface which will play the audio via your browser with the `--web` flag, the default port is `:8001` but this can be changed with the `--web-bind` option. Should work in any browser with support for Websockets and MP3 HTML5 audio.

If you're on OS/X, run

```
./tfl-announcements --afplay
```

to have the generated announcement played via the `afplay` command.

See `--help` for more flags.
