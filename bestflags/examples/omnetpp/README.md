# `omnetpp`

This directory contains the omnetpp example in SPEC2006 benchmark.

It also contains the json configuration file which includes the meta data
information to run the experiment.

This directory contains a build file `build_omnetpp` which is used by the build
module of the framework to compile the application.
This directory contains a test file `test_omnetpp` which is used by the test
module of the framework to benchmark the optimization compilation.
This directory contains a conf file which includes the set of optimization flags
the experiment will try.

To use this direction, first gives the file the executable permission.

```
chmod a+x build_bikjmp
chmod a+x test_bikjmp
```

Copy the SPEC2006 benchmark into this directory.

To run, invoke the `example_algorithm.py` in the parent directory.

```
python example_algorithms.py --file=examples/omnetpp/example.json
```

For help,

```
python example_algorithms.py --help
```
