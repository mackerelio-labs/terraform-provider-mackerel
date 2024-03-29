Contribution Guide
==================
We love contributions in the form of pull requests!
If you fixed or added something useful to the project, you can send pull-request.
Here's a quick guide:

1. Fork it from https://github.com/mackerelio-labs/terraform-provider-mackerel/fork
1. Create your feature branch (`git switch -c my-new-feature`)
1. Run test suite*1
1. If you add some new resource, please add documentation.
1. Commit your changes (`git commit -am 'Add some feature'`)
1. Push to the branch (`git push origin my-new-feature`)
1. Create new Pull Request*1

*1)  
PR needs to show that the changes passed the test in your local machine, so you have to paste the result of `$ make testacc TESTS=TestAccXXX`.  
Environment variables are required to run tests.  
`export MACKEREL_API_KEY=<YOUR-API-KEY>`  
Additional environment variables are required for AWS Integration.  
`export AWS_ROLE_ARN`, `export EXTERNAL_ID` or  
`export AWS_ACCESS_KEY_ID`, `export AWS_SECRET_ACCESS_KEY`  
You can run specific tests by giving a function name to `TESTS`.  
ex)
```zsh
$ make testacc TESTS=TestAccMackerelAWSIntegrationIAMRole    
TF_ACC=1 go test -v ./mackerel/... -run TestAccMackerelAWSIntegrationIAMRole -timeout 120m
=== RUN   TestAccMackerelAWSIntegrationIAMRole
=== PAUSE TestAccMackerelAWSIntegrationIAMRole
=== CONT  TestAccMackerelAWSIntegrationIAMRole
--- PASS: TestAccMackerelAWSIntegrationIAMRole (8.11s)
PASS
ok      github.com/mackerelio-labs/terraform-provider-mackerel/mackerel       8.701s
```
