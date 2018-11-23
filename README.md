# tc-client-generator

<img align="right" src="https://avatars3.githubusercontent.com/u/6257436?s=256" />Generate [taskcluster](https://tools.taskcluster.net/) clients in a variety of programming languages.

[![Taskcluster CI Status](https://github.taskcluster.net/v1/repository/taskcluster/tc-client-generator/master/badge.svg)](https://github.taskcluster.net/v1/repository/taskcluster/tc-client-generator/master/latest)
[![Linux Build Status](https://img.shields.io/travis/taskcluster/tc-client-generator.svg?style=flat-square&label=linux+build)](https://travis-ci.org/taskcluster/tc-client-generator)
[![GoDoc](https://godoc.org/github.com/taskcluster/tc-client-generator?status.svg)](https://godoc.org/github.com/taskcluster/tc-client-generator)
[![Coverage Status](https://coveralls.io/repos/taskcluster/tc-client-generator/badge.svg?branch=master&service=github)](https://coveralls.io/github/taskcluster/tc-client-generator?branch=master)
[![License](https://img.shields.io/badge/license-MPL%202.0-orange.svg)](http://mozilla.org/MPL/2.0)



# Install binary

* Download the latest release for your platform from https://github.com/taskcluster/tc-client-generator/releases
* For darwin/linux, make the binaries executable: `chmod a+x tc-client-generator*`

# Build from source

If you prefer not to use a prepackaged binary, or want to have the latest unreleased version from the development head:

* Head over to https://golang.org/dl/ and follow the instructions for your platform. __Note, go 1.11 or higher is required__.
* Run `go get github.com/taskcluster/tc-client-generator`

All being well, the binaries will be built under `$(go env GOPATH)/bin`.

# Acquire taskcluster credentials for running tests

There are two alternative mechanisms to acquire the scopes you need.

## Option 1

This method works if you log into Taskcluster via mozillians, *or* you log into
taskcluster via LDAP *using the same email address as your mozillians account*,
*or* if you do not currently have a mozillians account but would like to create
one.

* Sign up for a [Mozillians account](https://mozillians.org/en-US/) (if you do not already have one)
* Request membership of the [taskcluster-contributors](https://mozillians.org/en-US/group/taskcluster-contributors/) mozillians group

## Option 2

This method is for those who wish not to create a mozillians account, but
already authenticate into taskcluster via some other means, or have a
mozillians account but it is registered to a different email address than the
one they use to log into Taskcluster with (e.g. via LDAP integration).

* Request the scope `assume:project:taskcluster:tc-client-generator-tester` to be
  granted to you via a [bugzilla
  request](https://bugzilla.mozilla.org/enter_bug.cgi?product=Taskcluster&component=Service%20Request),
  including your [currently active `ClientId`](https://tools.taskcluster.net/credentials/)
  in the bug description. From the ClientId, we will be able to work out which role to assign the scope
  to, in order that you acquire the scope with the client you log into Taskcluster tools site with.

Once you have been granted the above scope:

* If you are signed into tools.taskcluster.net already, **sign out**
* Sign into [tools.taskcluster.net](https://tools.taskcluster.net/) using either your new Mozillians account, _or_ your LDAP account **if it uses the same email address as your Mozillians account**
* Check that a role or client of yours appears in [this list](https://tools.taskcluster.net/auth/scopes/assume%3Aproject%3Ataskcluster%3Atc-client-generator-tester)
* Create a permanent client (taskcluster credentials) for yourself in the [Client Manager](https://tools.taskcluster.net/auth/clients/) granting it the single scope `assume:project:taskcluster:tc-client-generator-tester`

To see a full description of all the config options available to you, run `tc-client-generator --help`:

```
tc-client-generator 0.0.1

tc-client-generator generates taskcluster clients in a variety of programming languages.

  Usage:
    tc-client-generator --help
    tc-client-generator --version

  Exit Codes:
    0      Completed successfully.
```

# Run the client generator

Simply run:

```
tc-client-generator
```

# Run the tc-client-generator test suite

For this you need to have the source files (you cannot run the tests from the binary package).

```
cd $(go env GOPATH)/src/github.com/taskcluster/tc-client-generator
./build.sh -t  ### add '-a' if you also want to build for all supported platforms
```


# Making a new tc-client-generator release

Run the `release.sh` script like so:

```
$ ./release.sh 0.0.1
```

This will perform some checks, tag the repo, push the tag to github, which will then trigger travis-ci to run tests, and publish the new release.

# Release notes

[tc-client-generator 0.0.1](https://github.com/taskcluster/tc-client-generator/releases/tag/v0.0.1)

This is an initial release, which doesn't do anything, but was made in order to test the release script was working.

# Further information

Please see:

* [Taskcluster Documentation](https://docs.taskcluster.net/)
