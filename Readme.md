



| Commands      | Sub Commands | Arguments                           | Flags  | Comments                      |
| ------------- | ------------ | ----------------------------------- | ------ | ----------------------------- |
| mailbox       | **add**      | user@mydomain.com plain_password    | --quota |                               |
|               | remove       | user@mydomain.com                   |        | delete aliases & forwardings? |
|               | **list**     |                                     | --filter       |  |
|  | add-alias | somealias user@mydomain.com |||
|  | remove-alias | alias@mydomain.com                  |||
|  | add-forwarding | erlec@bachmnaity.com |||
|  | remove-forwarding | erlec@bachmnaity.com |||
|  | **info**     | user@mydomain.com |||
| alias | **add**      | alias@mydomain.com user@mydomain.com |        |                               |
|               | **remove**   | alias@mydomain.com                  |        |                               |
|               | list   |                   | —filter |                               |
| forwarding    | add          | user@domain.com forward@example.com |        |                               |
|               | remove       | user@domain.com forward@example.com |||
|               | list |                                     | —filter |                               |
| catchall      | add          | domain.com dest@example.com |        |                               |
|               | remove       | domain.com dest@example.com |        |                               |
| domain | list |                                     |        |                               |
|  | info |                                     |        |                               |
|  | update |                                     |        |                               |

