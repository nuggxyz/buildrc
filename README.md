# buildrc 📦

`buildrc` is a command-line utility and GitHub Action designed to simplify and streamline your build processes. This tool was built with versatility in mind, allowing it to be easily integrated into various development workflows and CI/CD pipelines.

## `buildrc` Command

The `buildrc` command provides functionality to build and manage your projects. It checks for a "buildrc-override" artifact first; if it doesn't exist, `buildrc` will fetch a specific tagged version from GHCR (GitHub Container Registry), add the files to a temporary folder, and alias `buildrc` to the main executable.

### Usage

Here is a basic usage example:

```bash
$ buildrc <options>
```

## `buildrc` GitHub Action

The `buildrc` GitHub Action sets up `buildrc` in your GitHub Actions workflows. It downloads the "buildrc-override" artifact if it exists or pulls a specific version of `buildrc` from GHCR if the artifact does not exist.

The action then creates a temporary directory and sets an alias for `buildrc` in the current shell to point to the main executable.

### Usage

Here's a simple example of how to use the `buildrc` action in a workflow:

```yaml
- name: Setup buildrc
  uses: nuggxyz/actions/setup-buildrc@v0
```
