# ldapssh

`ldapssh` is a small Go script designed to be used by OpenSSH to get the public keys for a user stored in LDAP.
To use it, edit `/etc/ssh/sshd_config` and add the following lines:

```
AuthorizedKeysCommand /path/to/ldapssh
AuthorizedKeysCommandUser nslcd
```

Whatever user you use will require read access to `/etc/nslcd.conf` as that is where this program pulls the login
details for LDAP from.

Compile using Go.


## Configuration

This program will attempt to get all of the configuration it needs by reading the `nslcd` configuration at
`/etc/nslcd.conf`. If you do not use this, you will still have to create the configuration file.

```
uri ldap://ldap.mydomain.int/

base dc=mydomain,dc=int

filter passwd (objectClass=user)

binddn cn=admin,dc=mydomain,dc=int
bindpw mysupersecretpassword
```

The `uri` and `base` parameters are the only required ones. If you can't bind anonymously you will need to provide
`binddn` and `bindpw` and give it access to an account that can read the `sshPublicKey` attribute. If your users have
an `objectClass` other than `posixAccount`, you will also have to provide an appropriate filter. The filter provided
must be for the `passwd` service, like in the example above. Entries will be filtered with the filter
`(&<custom filter>(uid=username))` where `<custom filter>` is either `(objectClass=posixAccount)` or whatever you
specify in the `filter passwd` configuration option, and `username` is the username passed by OpenSSH.