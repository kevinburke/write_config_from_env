# write_config_from_env

This provides tools for pulling configuration values from environment variables,
and writing them to a YML file, or printing them to stdout as YAML.

If this is your environment:

```
FLAG_NEW_ROUTER=true
PORT=56789
PUBLIC_HOST=localhost:7000
TWILIO_ACCOUNT_SID=AC123
TIMEZONES=America/Los_Angeles,America/New_York
```

You will get a config file like this:

```
flag_new_router: "true"
port: "56789"
public_host: 'localhost:7000'
twilio_account_sid: AC123
timezones:
  - America/Los_Angeles
  - America/New_York
```

Note, all values are strings or arrays of strings. If a variable contains a
comma, it will be destructured into an array.

### Why not just use environment variables?

Environment variables leak trivially to exception reporting tools like Sentry,
or error pages/admin screens. Environment variables also get passed to
subprocesses by default. Putting your command in a Docker container does not
protect you from this.

Writing configuration to a file isn't perfect - people can still try to read the
file - but it's less error prone than using environment variables, and we set
the file permissions to 400, so the only reader is the current user.

Note, you have to remove the environment variables from the current process
before starting the command you care about. [The `envdir` program][envdir] is
useful for this. Create a series of empty files in a directory, where the
filenames are the vars you want to remove:

```
$ tree env
env
├── BASIC_AUTH_PASSWORD
├── BASIC_AUTH_USER
├── EMAIL_ADDRESS
├── ERROR_REPORTER_TOKEN
├── GOOGLE_CLIENT_ID
├── GOOGLE_CLIENT_SECRET
├── PUBLIC_HOST
├── SECRET_KEY
└── TWILIO_AUTH_TOKEN
```

Then run:

```
exec envdir ./env mycommand --config=config.yml
```

where `mycommand` is whatever command you want to run, and `--config=config.yml`
is your command loading configuration from the value you just loaded.

[envdir]: https://cr.yp.to/daemontools/envdir.html

### Usage

```
$ write_config_from_env -h
write_config_from_env

Read configuration from environment variables and write it to a yml file. By
default this script prints the config to stdout. Pass --config=<file> to write
to a file instead.

Usage of write_config_from_env:
  -config string
        Path to a config file
  -whitelist value
        Environment variables to whitelist. If unspecified, all environment
        variables will be written to the config
```

If you don't want to write every environment variable to the file, you can
specify a whitelist argument, as many times as you want.

```
write_config_from_env --whitelist=PATH --whitelist=PORT --whitelist=TZ
```

Only environment variables matching the whitelist will get written.

### Installation

Releases for popular operating systems are available on the [releases
page][releases]. To install them:

Find your target operating system (darwin, windows, linux) and desired bin
directory, and modify the command below as appropriate:

    curl --silent --location https://github.com/kevinburke/write_config_from_env/releases/download/0.5/write_config_from_env-linux-amd64 > /usr/local/bin/write_config_from_env && chmod 755 /usr/local/bin/write_config_from_env

On Travis, you may want to create `$HOME/bin` and write to that, since
/usr/local/bin isn't writable with their container-based infrastructure.


[releases]: https://github.com/kevinburke/write_config_from_env/releases
