# Strawberry - the main ingredient of your JAM stack

Strawberry is a Static Site Generator that is extremely fast, efficient, and modular.
It is written in [Go][go-site] and is a hard fork of [Hugo][hugo-site](v0.129.x).

[![Go Report Card](https://goreportcard.com/badge/github.com/strawberry-tools/strawberry)](https://goreportcard.com/report/github.com/strawberry-tools/strawberry)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/strawberry-tools/strawberry)
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/strawberry-tools/strawberry/tree/trunk.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/strawberry-tools/strawberry/tree/trunk)
[![codecov](https://codecov.io/gh/strawberry-tools/strawberry/branch/trunk/graph/badge.svg)](https://codecov.io/gh/strawberry-tools/strawberry)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](http://makeapullrequest.com)
[![GoDoc](https://godoc.org/github.com/strawberry-tools/strawberry?status.svg)](https://pkg.go.dev/github.com/strawberry-tools/strawberry)

**Strawberry is a relatively new project forked from Hugo.
Please stay tuned while we continue to update this readme and convert things over.**

## Overview

Strawberry is a static HTML and CSS website generator written in [Go][].
It is optimized for speed, ease of use, and configurability.
Strawberry takes a directory with content and templates and renders them into a full HTML website.

Strawberry relies on Markdown files with front matter for metadata, and you can run Strawberry from any directory.
This works well for shared hosts and other systems where you don’t have a privileged account.

Strawberry renders a typical website of moderate size in a fraction of a second.
A good rule of thumb is that each piece of content renders in around 1 millisecond.

Strawberry is designed to work well for any kind of website including docs, blogs, marketing sites, newsletters, and more.

#### Supported Architectures

Strawberry is built for Linux, macOS, and Windows for amd64 (Intel and AMD) CPUs.
Support for ARM is planned.

**Complete documentation is available at [Strawberry Documentation](https://gohugo.io/getting-started/).**

## Choose How to Install

If you want to use Strawberry as your site generator, simply install the Strawberry binary.
The Strawberry binary have no external dependencies.

To contribute to the Strawberry source code or documentation, you should [fork the Strawberry GitHub project](https://github.com/strawberryssg/strawberry-v0#fork-destination-box) and clone it to your local machine.

Finally, you can install the Strawberry source code with `go`, build the binary yourself, and run Strawberry that way.
Building the binary is an easy task for an experienced `go` getter.

### Install Strawberry as Your Site Generator (Binary Install)

Use the [installation instructions in the Strawberry documentation](https://gohugo.io/getting-started/installing/).

### Build and Install the Binaries from Source (Advanced Install)

#### Prerequisite Tools

* [Git](https://git-scm.com/)
* [Go (we test it with the last 2 major versions; but note that Hugo 0.95.0 only builds with >= Go 1.18.)](https://golang.org/dl/)

#### Fetch from GitHub

Strawberry uses the Go Modules support built into Go 1.11 to build. The easiest is to clone Strawberry in a directory outside of `GOPATH`, as in the following example:

```bash
mkdir $HOME/src
cd $HOME/src
git clone https://github.com/strawberry-tools/strawberry.git
cd strawberry
go install
```

**If you are a Windows user, substitute the `$HOME` environment variable above with `%USERPROFILE%`.**
	
## Contributing to Strawberry

For a complete guide to contributing to Strawberry, see the [Contribution Guide](CONTRIBUTING.md).

We welcome contributions to Strawberry of any kind including documentation, themes,
organization, tutorials, blog posts, bug reports, issues, feature requests,
feature implementations, pull requests, answering questions on the forum,
helping to manage issues, etc.

The Strawberry community and maintainers are [very active](https://github.com/strawberry-tools/strawberry/pulse/monthly) and helpful, and the project benefits greatly from this activity.

### Asking Support Questions

We're creating a forum soon where users and developers can ask questions.
You can use the GitHub Issue tracker to ask questions for now.

### Reporting Issues

If you believe you have found a defect in Strawberry or its documentation, use
the GitHub issue tracker to report the problem to the Strawberry maintainers.
When reporting the issue, please provide the version of Strawberry in use (`strawberry version`).

### Submitting Patches

The Strawberry project welcomes all contributors and contributions regardless of skill or experience level.
If you are interested in helping with the project, we will help you with your contribution.
Strawberry is a very active project with many contributions happening daily.

We want to create the best possible product for our users and the best contribution experience for our developers,
we have a set of guidelines which ensure that all contributions are acceptable.
The guidelines are not intended as a filter or barrier to participation.
If you are unfamiliar with the contribution process, the Strawberry team will help you and teach you how to bring your contribution in accordance with the guidelines.

For a complete guide to contributing code to Strawberry, see the [Contribution Guide](CONTRIBUTING.md).

[Go]: https://golang.org/
[Hugo Documentation]: https://gohugo.io/overview/introduction/
[go-site]: https://go.dev/
[hugo-site]: https://gohugo.io/
