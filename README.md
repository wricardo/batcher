batcher [![Build Status](https://travis-ci.org/wricardo/batcher.svg?branch=master)](https://travis-ci.org/wricardo/batcher) [![Coverage Status](https://coveralls.io/repos/wricardo/batcher/badge.png)](https://coveralls.io/r/wricardo/batcher)
=======

Batcher is a library that helps you to perform batch operations in Go. This library consists of "Collectors" and "Flushers". A Collector stores data in memory to be flushed latter. A Flusher receives data from a Collector and performs an operation in batch.
