Contribution Guide
==================
We love contributions in the form of pull requests!
If you fixed or added something useful to the project, you can send pull-request.
Here's a quick guide:

1. Fork it from https://github.com/mackerelio-labs/terraform-provider-mackerel/fork
1. Create your feature branch (`git switch -c my-new-feature`)
1. Run test suite with the `make test` or `go test ./...` command and confirm that it passes
   - This test suite requires environment variables, `MACKEREL_API_KEY`, `EXTERNAL_ID`, `AWS_ROLE_ARN`, `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.
1. If you add some new resource, please add documentation.
1. Commit your changes (`git commit -am 'Add some feature'`)
1. Push to the branch (`git push origin my-new-feature`)
1. Create new Pull Request
