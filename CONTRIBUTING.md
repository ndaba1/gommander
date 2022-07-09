# Contributions

All contributions are welcome and highly appreciated! There are a number of ways that you can do so:

- Improving the package and components documentation
- Helping us fix bugs
- Adding new features to the package
- Improving the code coverage
- Even filing isssues counts!

Another important way you can contribute to the package is adding a new example of the usage in the [examples](./examples/) directory.

# Goals of this package

The goals of this package are very simple. To provide a minimal and fast codebase for creating clis. As you make your contributions, please try to keep the no. of dependecies required by the package to a minimum. The only direct dependency of the package is the `fatih/color` pkg.
I'd also like to try and maintain backward compatibility but if your change proposes a breaking change, you can still submit it and we'll discuss it further.
Generally, the performance should be improved or at the least left the same after your proposed change. The package has scripts and make recipes for comparing various benchmarks.

# Testing

All the current tests are using custom-made assertion utils. Something like [testify]("https://github.com/u/testify") would also serve nicely but as mentioned, we're keeping the deps to a minimum. This may be changed in the future.
