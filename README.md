# cfg: The official Confbase CLI

[![Build Status](https://travis-ci.org/Confbase/cfg.svg?branch=master)](https://travis-ci.org/Confbase/cfg)


Hey! Confbase is still pre-alpha. It is not yet usable. Everything on this
page is subject to change at any point in time.

If you want to be an early adopter or if you want a free Confbase subscription,
send an email to thomas@confbase.com.

If you just want to know what the heck is going on, here you go:

Confbase...

* tries to make configuration management easier
* intends to solve problems similar to those solved by Ansible, Chef, and Puppet
    * intends to be used along side those products
    * intends to compete with those products in many ways
* leverages git
    * imposes a non-git-like workflow by default (see intended use section below)
    * allows users to use a traditional git workflow along side cfg

## Installation

Run `go get github.com/Confbase/cfg`

## Usage

After installation, the `cfg` binary will be in your `$GOBIN` directory.

Run `cfg -h` or `cfg --help` for usage.

## Primary features

### Templates and Instances

With `cfg`, you can mark a file as a "template" and another file as an "instance"
of that template. `cfg` then does a bunch of magic to help you not shoot yourself
in the foot.

For instance, you can't have a template with fields A, B, and C and
an instance with only fields A and B (no C), unless you explicitly say "C is an
optional field from now on." You also can't accidententally commit invalid JSON 
unless you explicitly say "it's okay for this file to be invalid JSON now."

Confbase also provides a bunch of tooling to help visualize and understand all
the different files in a base. E.g., "show me all my templates," "show me all the
instances of this template," "show me the difference between these instances"
and "tell me what's for dinner" are all wishes that can be fulfilled by `cfg`.

### Fetching config files on production machines

A raw copy of every file in the latest commit is made available via `cfg fetch`,
even if the machine running the command does not have git. These files are 
stored in an in-memory cache on Confbase servers for "blazing-fast" access.

`cfg snap new` makes the files from any commit permanently accessible via
`cfg fetch` (note that they can always be accessed through the underlying git
 repository). For example,

#### development box

```
$ cfg snap new my-new-feature
$ #... (update templates to account for new feature) ...
```

#### production box

```
$ # if you don't want to test the feature on live machines just yet

$ # this fetches config.yml from master
$ cfg fetch myteam:mybase/config.yml

$ # (or via curl and basic auth)
$ curl -u 'email:password' myteam.confbase.com/mybase/raw/master/config.yml
```

When the new feature is ready to be tested on live machines:

```
$ cfg fetch myteam:mybase/config.yml --snapshot my-new-feature
```

If all is well, the new feature will get pulled into master.


### Catching errors before commits

Pre-commit hooks are automatically added to the current base, based on its
contents.

For example, suppose a JSON template consisting of 12 fields is modified so that
a new field is added. Users will have to either add the new field to all
instances of the template or manually mark the field as optional. For another
example, if an instance of a JSON template is erroneously modified so that it is
no longer JSON, a pre-commit hook will automatically detect this
and ask the user to fix their error or use the `--force` flag.

`cfg lint` manually runs pre-commit hooks. The command also provides tools to
format config files in a uniform way.

### Inferring schemas

`cfg` provides tooling to generate schemas for config files based on example
data. In fact, it does this automatically when creating pre-commit hooks (which
it also does automatically).

## Intended workflow

cfg intends to be brutalist. All work is done on one branch. When the remote
branch gets ahead of the local branch, `cfg trampoline` or `cfg pull --hard`
are the two options provided by cfg to fix the predicament. If things get
really messy, `cfg git [...args]` provides direct access to the underlying
git repository.

### Running in an existing git repository

If a more traditional workflow is desired, consider running `cfg init --no-git`
and simply using git as before.
