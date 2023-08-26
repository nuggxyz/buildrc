# buildrc 📦
[![codecov](https://codecov.io/gh/walteh/buildrc/graph/badge.svg?token=6KW7C9X6YD)](https://codecov.io/gh/walteh/buildrc)
<!-- i have written a command line utility called "buildrc" that helps you version your git projects... can you help me write up a good readme?


basically you have a .buildrc in the root folder of your repository that contains a major version in this f9rmat { major: 3 } (it is flow yaml syntax) and then based on the latest tag/ref it will spit out a semver for you.


Right now there are three main situations it handles, commits to main, prs and local builds



basically when a commit to main happens, it will look at your full git history and find the latest semver tag in the chain. Then, it will either bump it by a minor version or a patch version based on the commit message. Baiscally if it has "patch" in it it will bump by patch.

when a pr happens ( -->


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


### Assumptions

1. You have a `.buildrc` file in the root of your repository.
2. pipeline will be run on linux-amd64
