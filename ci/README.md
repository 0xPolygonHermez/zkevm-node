# CI groups

In order to reduce the total time we spend in CI executions we are running the
end to end tests, lint and unit tests in parallel, so that the time of one
execution is roughly equal to the time spent on the longest parallel execution.

We have 3 different github actions workflows:
* `lint`, just runs the linter
* `test-full-non-e2e`, runs all non-e2e tests (unit, integration and functional)
* `test-e2e`, which uses a matrix strategy to run the e2e tests in 3 groups.

The e2e CI groups are defined in the `./ci/e2e-group{1,3}` directories. In each
directory we have symlinks that point to the actual e2e test to be executed (these
tests are defined under `./test/e2e`). The goal of these symlinks is keeping the
same code organization we have now while being able to run the costly e2e tests
in parallel on CI.

So, if for instance we have the following e2e tests defined:
* `./test/e2e/testA_test.go`
* `./test/e2e/testB_test.go`
* `./test/e2e/testC_test.go`
* `./test/e2e/testD_test.go`
and we want to run tests A and B in group1, test C in group 2 and test D in group 3
we would need to create these symlinks:
```
./ci/e2e-group1/testA_test.go -> ./test/e2e/testA_test.go
./ci/e2e-group1/testB_test.go -> ./test/e2e/testB_test.go
./ci/e2e-group2/testC_test.go -> ./test/e2e/testC_test.go
./ci/e2e-group3/testD_test.go -> ./test/e2e/testD_test.go
```
