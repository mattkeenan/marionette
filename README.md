# marionette

This is a proof of concept application which is designed to be puppet-like,
allowing you to define rules to get a system into a particular state.

At the moment there is support for:

* Defining rules.
  * Rules can have dependencies.
  * Rules can notify other rules when they're executed.
* Executing rules.

There are several builtin modules for performing tasks, and adding new primitives is simple thanks to the plugin-support documented later in this file.




# Installation & Usage

Install by executing:

    go get github.com/skx/marionette

Then launch the application with the path to a local set of rules, for example:

    marionette ./rules.txt



# Rule Definition

The general form of our rules looks like this:

```
   $MODULE [triggered] {
              name => "NAME OF RULE",
              arg_1 => "Value 1 ... ",
              arg_2 => [ "array values", "are fine" ],
              arg_3 => "Value 3 .. ",
   }
```

Each rule starts by declaring the type of module which is being invoked, then
there is a block containing "`key => value`" sections.  Different modules will
accept and expect different keys to configure themselves.  (Unknown arguments
will be ignored.)

In addition to the general arguments passed to the available modules you can also specify dependencies via two magical keys within each rule-block:

* `dependencies`
  * A list of any number of rules which must be executed before this.
* `notify`
  * A list of any number of rules which should be notified, because this rule was triggered.
    * Triggered in this sense means that the rule was executed and the state changed.

As a concrete example you need to run a command which depends upon a directory being present.  You could do this like so:


```
shell { name         => "Test",
        command      => "uptime > /tmp/blah/uptime",
        dependencies => [ "Create /tmp/blah" ] }

directory{ name   => "Create /tmp/blah",
           target => "/tmp/blah" }
```

The alternative would have been to have the directory-creation trigger the shell-execution rule via an explicit notification:

```
directory{ name   => "Create /tmp/blah",
           target => "/tmp/blah",
           notify => [ "Test" ]
}

shell triggered { name         => "Test",
                  command      => "uptime > /tmp/blah/uptime",
}
```

The difference in these two approaches is how often things run:

* In the first case we always run `uptime > /tmp/blah/uptime`
  * We just make sure that _before_ that the directory has been created.
* In the second case we run the command only once.
  * We run it only after the directory is created.
  * Because the directory-creation triggers the notification only when the rule changes (i.e. the directory goes from being "absent" to "present").


# Module Types

We have a small number of primitives at the moment implemented in 100% pure golang.

Adding new primitives as plugins is described later in this document.



## `directory`

The directory module allows you to create a directory, or change the permissions of one.

Example usage:

```
directory {  name    => "My home should have a binary directory",
             target  => "/home/steve/bin",
             mode    => "0755",
}
```

Valid parameters are:

* `target` is a mandatory parameter, and specifies the directory to be operated upon.
* `owner` - Username of the owner, e.g. "root".
* `group` - Groupname of the owner, e.g. "root".
* `mode` - The mode to set, e.g. "0755".
* `state` - Set the state of the directory.
  * `state => "absent"` remove it.
  * `state => "present"` create it (this is the default).



## `dpkg`

This module allows purging a package, or set of packages:

```
dpkg { name => "Remove stuff",
       package => ["vlc", "vlc-l10n"] }
```

Only the `package` key is required.

In the future we _might_ have an `apt` module for installing new packages.  We'll see.


## `file`

The file module allows a file to be created, from a local file, or via a remote HTTP-source.

Example usage:

```
file {  name       => "fetch file",
        target     => "/tmp/steve.txt",
        source_url => "https://steve.fi/",
}

file {  name     => "write file",
        target   => "/tmp/name.txt",
        content  => "My name is Steve",
}
```

`target` is a mandatory parameter, and specifies the file to be operated upon.

There are three ways a file can be created:

* `content` - Specify the content inline.
* `source_url` - The file contents are fetched from a remote URL.
* `source` - Content is copied from the existing path.

Other valid parameters are:

* `owner` - Username of the owner, e.g. "root".
* `group` - Groupname of the owner, e.g. "root".
* `mode` - The mode to set, e.g. "0755".
* `state` - Set the state of the file.
  * `state => "absent"` remove it.
  * `state => "present"` create it (this is the default).



## `git`

Clone a remote repository to a local directory.

Example usage:

```
git { path => "/tmp/xxx",
      repository => "https://github.com/src-d/go-git",
}
```

Valid parameters are:

* `repository` Contain the HTTP/HTTPS repository to clone.
* `path` - The location we'll clone to.
* `branch` - The branch to switch to, or be upon.
  * A missing branch will not be created.

If this module is used to `notify` another then it will trigger such a
notification if either:

* The repository wasn't present, and had to be cloned.
* The repository was updated.  (i.e. Remote changes were pulled in.)


## `link`

The `link` module allows you to create a symbolic link.

Example usage:

```
link { name => "Symlink test",
       source => "/etc/passwd",  target => "/tmp/password.txt" }
```

Valid parameters are:

* `target` is a mandatory parameter, and specifies the location of the symlink to create.
* `source` is a mandatory parameter, and specifies the item the symlink should point to.


## `shell`

The shell module allows you to run shell-commands.

Example:

```
shell { name => "I touch your file.",
        command => "touch /tmp/blah/test.me" }
```

`command` is the only mandatory parameter.



# Adding Modules

Adding marionette modules can be done in two ways:

* Implementing your plugin in 100% go.
  * As you'll see done for the existing [modules/](modules/) we include.
* Writing a plugin.
  * You can write a plugin in __any__ language you like.

Given a rule such as the following we'll look for a handler for `foo`:

```
foo { name => "My rule",
      param1 => "One", }
```

If there is no built-in plugin with that name then the command located at `~/.marionette/foo` will be executed.  The parameters will be passed to that binary by being as JSON piped to STDIN.

If that plugin makes a change, such that triggers should be executed, it should print `changed` to STDOUT and exit with a return code of 0.  If no change was made then it should print `unchanged` to STDOUT and also exit with a return code of 0.

A non-zero return code will be assumed to mean something failed, and execution will terminate.



# Example

See [input.txt](input.txt) for a sample rule-file, including syntax breakdown.
