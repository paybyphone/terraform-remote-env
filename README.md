terraform-remote-env
=====================

This is a simple remote state helper for [Terraform][1] - it will pull all
outputs from a remote state and print them out as `TF_VAR_foo` environment
variables.

## Why?

Current limitations in Terraform do not allow the easy use of resource
variables, including ones obtained from `terraform_remote_state`, in things like
provider configuration, or a resource's `count`.

This is something that is currently being worked on, see [here][2] and
[here][3], however these features may be a bit of a ways away. This tool is
meant to be a holdover, and may be deprecated when this comes, unless you have
need for remote state outputs for something else, of course.

## How?

You can build this tool yourself using `go get -u
github.com/paybyphone/terraform-remote-env`, or head over to the [releases][4]
page to get the latest binary release.

## Options

This tool essentially takes the same options as `terraform remote config` does
see [here][5].

Output is in a single-line  `TF_VAR_foo=bar TF_VAR_baz=qux` format. With the
`-prefix` flag, you can specify a prefix to add to the `TF_VAR_` variable, ie:
`TF_VAR_prefix_foo=bar TF_VAR_prefix_baz=qux`.

[1]: https://terraform.io/
[2]: https://github.com/hashicorp/terraform/pull/6598
[3]: https://github.com/hashicorp/terraform/issues/4149
[4]: https://github.com/paybyphone/terraform-remote-env/releases
[5]: https://www.terraform.io/docs/commands/remote-config.html

## License 

```
Copyright 2016 PayByPhone Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
