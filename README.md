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
