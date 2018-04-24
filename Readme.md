



| Commands      | Sub Commands | Arguments                           | Flags  | Comments                      |
| ------------- | ------------ | ----------------------------------- | ------ | ----------------------------- |
| mailbox       | add          | user@mydomain.com plain_password    | --quota |                               |
|               | remove       | user@mydomain.com                   |        | delete aliases & forwardings? |
|               | list         |                                     | --filter       |  |
| alias-account | add          | alias@mydomain.com                  |        |                               |
|               | remove       | alias@mydomain.com                  |        |                               |
| forwarding    | add          | user@domain.com forward@example.com |        |                               |
|               | remove       | user@domain.com forward@example.com |||
|               |        |                                     |        |                               |
| catchall      | add          | domain.com dest@example.com |        |                               |
|               | remove       | domain.com dest@example.com |        |                               |
|               |              |                                     |        |                               |

