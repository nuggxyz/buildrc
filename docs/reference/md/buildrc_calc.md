## buildrc calc

calculate next pre-release tag

```
buildrc calc [flags]
```

### Options

```
  -c, --commit-message-override string   The commit message to use
  -h, --help                             help for calc
  -l, --latest-tag-override string       The tag to use
  -m, --major uint                       The major version to set
  -p, --patch                            shortcut for --type=release --patch-indicator=x --commit-message-override=x
  -i, --patch-indicator string           The ref to calculate the patch from (default "patch")
  -n, --pr-number int                    The pr number to set
  -t, --type string                      The type of commit to calculate (default "local")
```

### Options inherited from parent commands

```
  -d, --debug            Print debug output
  -f, --file string      The buildrc file to use (default ".buildrc")
  -g, --git-dir string   The git directory to use (default ".")
  -q, --quiet            Do not print any output
  -v, --version          Print version and exit
```

### SEE ALSO

* [buildrc](buildrc.md)	 - buildrc brings tabs to terraform

