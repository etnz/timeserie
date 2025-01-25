# timeserie

'time' package based time series.

This package deals with timeseries. 

It defines a Support, that holds couples of time.Time and float64 value in chronological order.

It defines a Function, based on a support, that can support different operations (add, mult, etc.).

It provides a jsonline format to serialize Supports into a value change dump format.

It provides utilities function to deal with filtering, grouping, sampling timeserie Supports.