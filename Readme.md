# iredmail-cli

[![CircleCI](https://circleci.com/gh/drlogout/iredmail-cli/tree/master.svg?style=svg)](https://circleci.com/gh/drlogout/iredmail-cli/tree/master)

## Commands

###mailbox

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
#### info

Show mailbox info.

Example:

```bash
$ iredmail-cli mailbox info info@example.com

+----------------------+---------------------------------------------+
|       MAILBOX        |              info@example.com               |
+----------------------+---------------------------------------------+
| Quota                | 2048 MB                                     |
| Forwardings          | info@otherdomain.com                        |
|                      | webmaster@otherexample.net                  |
| Keep copy in mailbox | yes                                         |
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

#### update

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

#### add-alias

Add a mailbox alias.

A mailbox `info@example.com` can have additional email addresses like `abuse@example.com`, `webmaster@example.com` and more, all emails sent to these addresses will be delivered to same mailbox (`info@example.com`). 

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
|                      | webmaster									 |
| Forwardings          | info@otherdomain.com                        |
|                      | webmaster@otherexample.net                  |
| Keep copy in mailbox | yes                                         |
| Maildir              | example.com/i/n/f/info-2018.07.09.09.13.27/ |
+----------------------+---------------------------------------------+
```

#### delete-alias

Delete an alias.

Example:

```bash
$ iredmail-cli mailbox delete-alias abuse@example.com
```

### forwarding

#### add



#### delete

#### list