# iredmail-cli

[![CircleCI](https://circleci.com/gh/drlogout/iredmail-cli/tree/master.svg?style=svg)](https://circleci.com/gh/drlogout/iredmail-cli/tree/master)

## Table of contents

* [Commands](#commands)
    * [mailbox](#mailbox)
        * [add [MAILBOX_EMAIL] [PLAIN_PASSWORD]](#add-mailbox_email-plain_password)
        * [delete [MAILBOX_EMAIL]](#delete-mailbox_email)
        * [info [MAILBOX_EMAIL]](#info-mailbox_email)
        * [list](#list)
        * [update [MAILBOX_EMAIL]](#update-mailbox_email)
        * [add-alias [ALIAS] [MAILBOX_EMAIL]](#add-alias-alias-mailbox_email)
        * [delete-alias [ALIAS_EMAIL]](#delete-alias-alias_email)
    * [forwarding](#forwarding)
        * [add [MAILBOX_EMAIL] [DESTINATION_EMAIL]](#add-mailbox_email-destination_email)
        * [delete [MAILBOX_EMAIL] [DESTINATION_EMAIL]](#delete-mailbox_email-destination_email)
        * [list](#list-1)
    * [domain](#domain)
        * [add [DOMAIN]](#add-domain)
        * [delete [DOMAIN]](#delete-domain)
        * [list](#list-2)
        * [add-alias [ALIAS_DOMAIN] [DOMAIN]](#add-alias-alias_domain-domain)
        * [delete-alias  [ALIAS_DOMAIN]](#delete-alias--alias_domain)
        * [add-catchall [DOMAIN] [DESTINATION_EMAIL]](#add-catchall-domain-destination_email)
        * [delete-catchall [DOMAIN] [DESTINATION_EMAIL]](#delete-catchall-domain-destination_email)
    * [alias](#alias)
        * [add [ALIAS_EMAIL]](#add-alias_email)
        * [delete [ALIAS_EMAIL]](#delete-alias_email)
        * [info [ALIAS_EMAIL]](#info-alias_email)
        * [list](#list-3)
        * [add-forwarding [ALIAS_EMAIL] [DESTINATION_EMAIL]](#add-forwarding-alias_email-destination_email)
        * [delete-forwarding [ALIAS_EMAIL] [DESTINATION_EMAIL]](#delete-forwarding-alias_email-destination_email)
    * [version](#version)

## Commands

To print the help of a command or sub command append the `—help` or `-h` flag.

------

### mailbox

Add/delete/list mailboxes and mailbox-aliases.

#### add \[MAILBOX_EMAIL] \[PLAIN_PASSWORD]

Add a new mailbox.

Example:

```bash
$ iredmail-cli mailbox add info@example.com swekjhlwekjdhw
```

Flags:

- -q, --quota: Set custom quota in MB, default 2048
- -s, --storage-path: Set custom storage path

#### delete \[MAILBOX_EMAIL]

Delete a mailbox.

Example:

```bash
$ iredmail-cli mailbox delete info@example.com
```
Flags:

- -f, --force: Force deletion

#### info \[MAILBOX_EMAIL]

Show mailbox info.

Example:

```bash
$ iredmail-cli mailbox info info@example.com
+----------------------+---------------------------------------------+
|       MAILBOX        |              info@example.com               |
+----------------------+---------------------------------------------+
| Quota                | 2048 MB                                     |
| Maildir              | example.com/i/n/f/info-2018.07.09.09.13.27/ |
+----------------------+---------------------------------------------+
```
#### list

List mailboxes.

Example:

```bash
$ iredmail-cli mailbox list
+-----------------------+------------+
|        MAILBOX        | QUOTA (MB) |
+-----------------------+------------+
| info@domain.com       |       2048 |
| info@example.com      |       2048 |
| mail@example.net      |       2048 |
| support@example.com   |       2048 |
+-----------------------+------------+

# To filter results use the --filter flag
$ iredmail-cli mailbox list -f example.com
+-----------------------+------------+
|        MAILBOX        | QUOTA (MB) |
+-----------------------+------------+
| info@example.com      |       2048 |
| support@example.com   |       2048 |
+-----------------------+------------+
```

Flags:

- -f, --filter: Filter results

#### update \[MAILBOX_EMAIL]

Update keep-copy and quota.

If mailboxes with forwardings should not keep a copy of the forwarded email use "--keep-copy no".
This is only possible if at least one forwarding for [MAILBOX_EMAIL] exists.
By default copies are kept in the mailbox.

The quota of the mailbox could be set with this flag, e.g. "--quota 4096" (in MB).

Example:

```bash
$ iredmail-cli mailbox update info@example.com -k no
$ iredmail-cli mailbox update info@example.com -q 4098
```
Flags:

- -k, --keep-copy: enable or disable keep-copy
- -q, --quota: Set custom quota in MB

#### add-alias [ALIAS] \[MAILBOX_EMAIL]

Add a mailbox alias.

A mailbox `info@example.com` can have additional email addresses like `abuse@example.com`, `webmaster@example.com` and more, all emails sent to these addresses will be delivered to the same mailbox (`info@example.com`). 

Example:

```bash
$ iredmail-cli mailbox add-alias abuse info@example.com
$ iredmail-cli mailbox add-alias webmaster info@example.com
$ iredmail-cli mailbox info info@example.com
+----------------------+---------------------------------------------+
|       MAILBOX        |              info@example.com               |
+----------------------+---------------------------------------------+
| Quota                | 2048 MB                                     |
| Mailbox aliases      | abuse                                       |
|                      | webmaster                                   |
| Maildir              | example.com/i/n/f/info-2018.07.09.09.13.27/ |
+----------------------+---------------------------------------------+
```

#### delete-alias [ALIAS_EMAIL]

Delete an alias.

Example:

```bash
$ iredmail-cli mailbox delete-alias abuse@example.com
```
------

### forwarding

Add/delete/list forwardings.

#### add \[MAILBOX_EMAIL] [DESTINATION_EMAIL]

Add forwarding.

Example:

```bash
$ iredmail-cli forwarding add info@example.com post@otherdomain.com
$ iredmail-cli forwarding add info@example.com tech@company.com
$ iredmail-cli mailbox info info@example.com
+----------------------+---------------------------------------------+
|       MAILBOX        |              info@example.com               |
+----------------------+---------------------------------------------+
| Quota                | 2048 MB                                     |
| Mailbox aliases      | abuse                                       |
|                      | webmaster                                   |
| Forwardings          | tech@company.com                            |
|                      | post@otherdomain.com                        |
| Keep copy in mailbox | yes                                         |
| Maildir              | example.com/i/n/f/info-2018.07.09.09.13.27/ |
+----------------------+---------------------------------------------+

```

By default a copy will be left in the mailbox, to change that behavior use the `iredmail-cli mailbox update` command.

#### delete \[MAILBOX_EMAIL] [DESTINATION_EMAIL]

Delete forwarding.

Example:

```bash
$ iredmail-cli forwarding delete info@example.com tech@company.com
```

#### list

List forwardings.

Example:

```bash
$ iredmail-cli forwarding list
+------------------+----------------------------+----------------------+
|  MAILBOX EMAIL   |     DESTINATION EMAIL      | KEEP COPY IN MAILBOX |
+------------------+----------------------------+----------------------+
| info@example.com | tech@company.com           | yes                  |
|                  | post@otherdomain.com       |                      |
| mail@example.net | mail@domain.com            | no                   |
+------------------+----------------------------+----------------------+
```

Flags:

- -f, --filter: Filter results

------

### domain

Add/delete/list domains, domain aliases and catchall forwardings.

#### add [DOMAIN]

Add a domain.

Example:

```bash
$ iredmail-cli domain add somedomain.com
```

Flags:

- -d, --description: Domain description
- -s, --settings: Domain settings (default: default_user_quota:2048)

#### delete [DOMAIN]

Delete a domain.

Example:

```bash
$ iredmail-cli domain delete somedomain.com
```

Flags:

- -f, --force: Force deletion

#### list

List domains.

Example:

```bash
$ iredmail-cli domain list
+-------------+-----------+-------------------+-------------+
|   DOMAIN    |   ALIAS   | CATCH-ALL ADDRESS | DESCRIPTION |
+-------------+-----------+-------------------+-------------+
| domain.com  |           |                   |             |
| example.com |           |                   |             |
+-------------+-----------+-------------------+-------------+
```

Flags:

- -f, --filter: Filter results

#### add-alias \[ALIAS_DOMAIN] \[DOMAIN]

Add an alias domain.
Emails sent to user@[ALIAS_DOMAIN] will be delivered to user@[DOMAIN].

Example:

```bash
$ iredmail-cli domain add-alias domain.net domain.com
$ iredmail-cli domain list
+-------------+------------+-------------------+-------------+
|   DOMAIN    |   ALIAS    | CATCH-ALL ADDRESS | DESCRIPTION |
+-------------+------------+-------------------+-------------+
| domain.com  | domain.net |                   |             |
| example.com |            |                   |             |
+-------------+------------+-------------------+-------------+
```

#### delete-alias  \[ALIAS_DOMAIN]

Delete an alias domain.

Example:

```bash
$ iredmail-cli domain delete-alias [ALIAS_DOMAIN]
```

#### add-catchall \[DOMAIN] \[DESTINATION_EMAIL]

Add a per-domain catch-all forwarding.
Emails sent to non-existing mailboxes of [DOMAIN] will be delivered to [DESTINATION_EMAIL].
Multiple [DESTINATION_EMAIL]s are possible.

Example:

```bash
$ iredmail-cli domain add-catchall example.com info@example.com
$ iredmail-cli domain add-catchall example.com post@otherdomain.com

$ iredmail-cli domain list
+-------------+------------+----------------------+-------------+
|   DOMAIN    |   ALIAS    | CATCH-ALL ADDRESS    | DESCRIPTION |
+-------------+------------+----------------------+-------------+
| domain.com  | domain.net |                      |             |
| example.com |            | info@example.com     |             |
|             |            | post@otherdomain.com |             |
+-------------+------------+----------------------+-------------+
```

#### delete-catchall \[DOMAIN] \[DESTINATION_EMAIL]

Delete a per-domain catch-all forwarding.

Example:

```bash
$ iredmail-cli domain delete-catchall example.com post@otherdomain.com
```

------

### alias

Add/delete/list aliases and their forwardings.

#### add [ALIAS_EMAIL]

Add an alias.

Example:

```bash
$ iredmail-cli alias add tech@example.com
```

#### delete [ALIAS_EMAIL]

Delete an alias.

Example:

```bash
$ iredmail-cli alias delete tech@example.com
```

Flags:

- -f, --force: Force deletion

#### info [ALIAS_EMAIL]

Show alias info.

Example:

```bash
$ iredmail-cli alias info tech@example.com
+--------------------+---------------------------+
|       ALIAS        |        FORWARDINGS        |
+--------------------+---------------------------+
| tech@example.com   | info@example.com          |
|                    | chris@example.com         |
|                    | pete@domain.com           |
+--------------------+---------------------------+
```

#### list

List aliases.

Example:

```bash
$ iredmail-cli alias list
+-----------------------+---------------------------+
|         ALIAS         |        FORWARDINGS        |
+-----------------------+---------------------------+
| tech@example.com      |                           |
| help@example.net      |                           |
+-----------------------+---------------------------+
```

Flags:

- -f, --filter: Filter results

#### add-forwarding \[ALIAS_EMAIL] \[DESTINATION_EMAIL] 

Add forwarding to an alias.
Emails sent to [ALIAS_EMAIL] will be delivered to [DESTINATION_EMAIL].
An alias can have multiple forwardings.

Example:

```bash
$ iredmail-cli alias add tech@example.com info@exmaple.com
$ iredmail-cli alias add tech@example.com pete@domain.com

+-----------------------+---------------------------+
|         ALIAS         |        FORWARDINGS        |
+-----------------------+---------------------------+
| tech@example.com      | info@exmaple.com          |
|                       | pete@domain.com           |
+-----------------------+---------------------------+
```

#### delete-forwarding \[ALIAS_EMAIL] \[DESTINATION_EMAIL]

Delete forwarding from an alias.

Example: 

```bash
$ iredmail-cli alias delete tech@example.com pete@domain.com 
```

------

### version

Show iredMail and iredmail-cli version.
