batcher [![Build Status](https://travis-ci.org/wricardo/batcher.svg?branch=master)](https://travis-ci.org/wricardo/batcher) [![Coverage Status](https://coveralls.io/repos/wricardo/batcher/badge.png)](https://coveralls.io/r/wricardo/batcher)
=======

[![Join the chat at https://gitter.im/wricardo/batcher](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/wricardo/batcher?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Batcher is a library that helps you to perform batch operations in Go. This library consists of "Collectors" and "Flushers". A Collector stores data in memory to be flushed latter. A Flusher receives data from a Collector and performs an operation in batch.
