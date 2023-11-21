# CI groups

In order to reduce the total time we spend in CI executions we are running the
end to end tests, lint and unit tests in parallel, so that the time of one
execution is roughly equal to the time spent on the longest parallel execution.

We have 3 different github actions workflows:
* `lint`, just runs the linter
* `test-full-non-e2e`, runs all non-e2e tests (unit, integration and functional)
* `test-e2e`, which uses a matrix strategy to run the e2e tests, currently using
3 groups.

The e2e CI groups are defined in the `./ci/e2e-group{1,N}` directories. In each
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
## How to enable/disable groups
As stated above, the `test-e2e` workflow relies on a matrix strategy for executing
the tests in the different groups, both for the `trusted` and the `from-fork` jobs
(you can read more about these jobs [here](./ok-to-test.md)). The matrix strategy
for each of the jobs looks like this:
```
strategy:
  matrix:
    go-version: [ 1.21.x ]
    goarch: [ "amd64" ]
    e2e-group: [ 1, 2, 3 ]
```
and the step that executes the tests relies on make targets of the form `test-e2e-group-n`
and looks like this:
```
- name: Test
  run: make test-e2e-group-${{ matrix.e2e-group }}
  working-directory: test
```
If you want to disable a group, we just need to remove it from the `e2e-group`
array in the matrix strategy. Given the configuration above, if we want to disable
groups 1 and 3, the matrix strategy config should look like:
```
strategy:
  matrix:
    go-version: [ 1.21.x ]
    goarch: [ "amd64" ]
    e2e-group: [ 2 ]
```
If we want to re-add group 1:
```
strategy:
  matrix:
    go-version: [ 1.21.x ]
    goarch: [ "amd64" ]
    e2e-group: [ 1, 2 ]
```
## Add new groups
In order to add a new group, for instance, group 4, we should:
* create a new subdir `e2e-group4` under `ci/`.
* symlink from `./ci/e2e-group4` the e2e test (under `./test/e2e`) that we want
to belong to the new group.
* create a new makefile entry, `test-e2e-group-4` that executes the go tests in
`./ci/e2e-group4`, for instance:
```
.PHONY: test-e2e-group-4
test-e2e-group-4: build-docker compile-scs ## Runs group 4 e2e tests checking race conditions
	$(STOPDB)
	$(RUNDB); sleep 7
	trap '$(STOPDB)' EXIT; MallocNanoZone=0 go test -race -p 1 -timeout 600s ./ci/e2e-group4/...
```
* include `4` in the matrix definition for the `trusted` and `from-fork` jobs
in `.github/workflows/test-e2e.yml`

*NOTE*: Usually groups should be as packed as possible so that we can optimize
the number of test lanes and the total execution time. If, for instance, we have
a group with one single test that takes 10min we should try to add tests to the
other groups with a total execution time (adding the execution time of each test
in the group) with up to 10min.
