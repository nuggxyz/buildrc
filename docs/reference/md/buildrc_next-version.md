## buildrc next-version

calculate next pre-release tag

```
buildrc next-version [flags]
```

### Options

```
  -a, --auto                             shortcut for if CI != 'true' then local else if '--pr-number' > 0 then pr
  -c, --commit-message-override string   The commit message to use
      --git-dir string                   Git directory
  -h, --help                             help for next-version
  -l, --latest-tag-override string       The tag to use
      --no-v                             do not prefix with 'v'
  -p, --patch                            shortcut for --patch-indicator=x --commit-message-override=x
  -i, --patch-indicator string           The ref to calculate the patch from (default "patch")
  -n, --pr-number uint                   The pr number to set
  -t, --type string                      The type of commit to calculate (default "local")
```

### Options inherited from parent commands

```
  -d, --debug     Print debug output
  -q, --quiet     Do not print any output
  -v, --version   Print version and exit
```

### SEE ALSO

* [buildrc](buildrc.md)	 - build time metadata

