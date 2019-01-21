
|![](https://upload.wikimedia.org/wikipedia/commons/thumb/4/4c/Anchor_pictogram_yellow.svg/156px-Anchor_pictogram_yellow.svg.png) | Drycc Workflow is the open source fork of Drycc Workflow.<br />Please [go here](https://www.drycc.com/) for more detail. |
|---:|---|
| 08/27/2018 | Team Drycc [blog][] comes online |
| 08/20/2018 | Drycc [#community slack][] goes dark |
| 08/10/2018 | Drycc Workflow [v2.19.4][] fourth patch release |
| 08/08/2018 | [Drycc website][] goes dark, then redirects to Azure Kubernetes Service |
| 08/01/2018 | Drycc Workflow [v2.19.3][] third patch release |
| 07/17/2018 | Drycc Workflow [v2.19.2][] second patch release |
| 07/12/2018 | Drycc Workflow [v2.19.1][] first patch release |
| 06/29/2018 | Drycc Workflow [v2.19.0][] first release in the open source fork of Drycc |
| 06/16/2018 | Drycc Workflow [v2.19][] series is announced |
| 03/01/2018 | End of Drycc Workflow maintenance: critical patches no longer merged |
| 12/11/2017 | Team Drycc [slack community][] invites first volunteers |
| 09/07/2017 | Drycc Workflow [v2.18][] final release before entering maintenance mode |
| 09/06/2017 | Team Drycc [slack community][] comes online |

# Controller Go SDK
[![Build Status](https://ci.drycc.cc/buildStatus/icon?job=Drycc/controller-sdk-go/master)](https://ci.drycc.cc/job/Drycc/job/controller-sdk-go/job/master/)
[![codecov](https://codecov.io/gh/drycc/controller-sdk-go/branch/master/graph/badge.svg)](https://codecov.io/gh/drycc/controller-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/drycc/controller-sdk-go)](https://goreportcard.com/report/github.com/drycc/controller-sdk-go)
[![codebeat badge](https://codebeat.co/badges/2fdee091-714d-4860-ab19-dba7587a3158)](https://codebeat.co/projects/github-com-drycc-controller-sdk-go)
[![GoDoc](https://godoc.org/github.com/drycc/controller-sdk-go?status.svg)](https://godoc.org/github.com/drycc/controller-sdk-go)

This is the Go SDK for interacting with the [Drycc Controller](https://github.com/drycc/controller).

### Usage

```go
import drycc "github.com/drycc/controller-sdk-go"
import "github.com/drycc/controller-sdk-go/apps"
```

Construct a drycc client to interact with the controller API. Then, get the first 100 apps the user has access to.

```go
//                    Verify SSL, Controller URL, API Token
client, err := drycc.New(true, "drycc.test.io", "abc123")
if err != nil {
    log.Fatal(err)
}
apps, _, err := apps.List(client, 100)
if err != nil {
    log.Fatal(err)
}
```

### Authentication

```go
import drycc "github.com/drycc/controller-sdk-go"
import "github.com/drycc/controller-sdk-go/auth"
```

If you don't already have a token for a user, you can retrieve one with a username and password.

```go
// Create a client with a blank token to pass to login.
client, err := drycc.New(true, "drycc.test.io", "")
if err != nil {
    log.Fatal(err)
}
token, err := auth.Login(client, "user", "password")
if err != nil {
    log.Fatal(err)
}
// Set the client to use the retrieved token
client.Token = token
```

For a complete usage guide to the SDK, see [full package documentation](https://godoc.org/github.com/drycc/controller-sdk-go).

[v2.18]: https://github.com/drycc/workflow/releases/tag/v2.18.0
[k8s-home]: http://kubernetes.io
[install-k8s]: http://kubernetes.io/gettingstarted/
[mkdocs]: http://www.mkdocs.org/
[issues]: https://github.com/drycc/workflow/issues
[prs]: https://github.com/drycc/workflow/pulls
[Drycc website]: http://drycc.com/
[blog]: https://blog.drycc.info/blog/
[#community slack]: https://slack.drycc.cc/
[slack community]: https://slack.drycc.com/
[v2.18]: https://github.com/drycc/workflow/releases/tag/v2.18.0
[v2.19]: https://web.drycc.com
[v2.19.0]: https://gist.github.com/Cryptophobia/24c204583b18b9fc74c629fb2b62dfa3/revisions
[v2.19.1]: https://github.com/drycc/workflow/releases/tag/v2.19.1
[v2.19.2]: https://github.com/drycc/workflow/releases/tag/v2.19.2
[v2.19.3]: https://github.com/drycc/workflow/releases/tag/v2.19.3
[v2.19.4]: https://github.com/drycc/workflow/releases/tag/v2.19.4
