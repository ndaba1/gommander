# Contributions

All contributions are welcome and highly appreciated! There are a number of ways that you can do so:

- Improving the package and components documentation
- Helping us fix bugs
- Adding new features to the package
- Improving the code coverage
- Even filing isssues counts!

Another important way you can contribute to the package is adding a new example of the usage in the [examples](./examples/) directory.

# Goals of this package

The goals of this package are very simple:

1. To provide a minimal and fast codebase for creating clis. As you make your contributions, please try to keep the no. of dependecies required by the package to a minimum. The only direct dependency of the package is the `fatih/color` pkg.

2. I'd also like to try and maintain backward compatibility but if your change proposes a breaking change, you can still submit it and we'll discuss it further.

3. Generally, the performance should be improved or at the least left the same after your proposed change. The package has scripts and make recipes for comparing various benchmarks.

# Testing

All the current tests are using custom-made assertion utils. Something like [testify]("https://github.com/u/testify") would also serve nicely but as mentioned, we're trying to keep the dependencies to a minimum. This may be changed in the future.

The code coverage is also tracked with codecov and changes increasing the coverage are always welcome. If you make changes but don't write tests for them, these will cause the codecov status checks in the coverage workflow to fail upon your pull request. However, the PR can still be merged as tests can always be added later.

# Benches

Performance is very important to the package. We try to offer an ergonomic but fast package. You can run benches for the package by either using the make target for this `make bench` or using standard go tools.

Once you have cloned the repo, I'd advice that you first run a benchmark and echo the output to a file. This will be useful as it will be later used to compare the performance after you've made your changes. The easiest way to do this is to run the `./benches.sh` bash script.

The bash script runs the benches and echo's the output to a file who's name corresponds to the Unix Nano time i.e `date +"%s"`. It then symlinks the latest benchmark to a file named `latest.bench` and the previous one, if any, to `old.bench`.

When you've made your changes and are ready to compare them, run the script again and once its done, run the `make benchcmp` target to compare them.
