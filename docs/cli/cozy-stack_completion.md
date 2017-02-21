## cozy-stack completion

Output shell completion code for the specified shell

### Synopsis



Output shell completion code for the specified shell (bash or zsh).
The shell code must be evalutated to provide interactive
completion of kubectl commands.  This can be done by sourcing it from
the .bash_profile.

Note: this requires the bash-completion framework, which is not installed
by default on Mac.  This can be installed by using homebrew:

    $ brew install bash-completion

Once installed, bash_completion must be evaluated.  This can be done by adding the
following line to the .bash_profile

    $ source $(brew --prefix)/etc/bash_completion

```
cozy-stack completion [shell]
```

### Examples

```
# cozy-stack completion bash > /etc/bash_completion.d/cozy-stack
```

### Options inherited from parent commands

```
      --admin-host string      administration server host (default "localhost")
      --admin-port int         administration server port (default 6060)
      --assets string          path to the directory with the assets (use the packed assets by default)
  -c, --config string          configuration file (default "$HOME/.cozy.yaml")
      --couchdb-url string     CouchDB URL (default "http://localhost:5984/")
      --fs-url string          filesystem url (default "file://localhost//storage")
      --host string            server host (default "localhost")
      --log-level string       define the log level (default "info")
      --mail-disable-tls       disable smtp over tls
      --mail-host string       mail smtp host (default "localhost")
      --mail-password string   mail smtp password
      --mail-port int          mail smtp port (default 465)
      --mail-username string   mail smtp username
  -p, --port int               server port (default 8080)
      --subdomains string      how to structure the subdomains for apps (can be nested or flat) (default "nested")
```

### SEE ALSO
* [cozy-stack](cozy-stack.md)	 - cozy-stack is the main command
