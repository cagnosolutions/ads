# This repository, 'ads',  is an ambiguous acronym standing for either 'algorithms and data structures' or 'advanced data structures' 
---
It contains data structures, algorithms, and helpers for the go programming language. The items in this package are things that I have come across and needed to use at some time or another that are not contained within the standard library (as of go 1.5.) Most if not all of them have been used in some type of production environment. Most of them should not be considered thread-safe, they will note if they are.

* Currently ads contains the following:

> 'ads/bpt' is a lexigraphically ordered b+tree using []byte for both the internal keys as well as the leaf node values.
> 'ads/mio' is a fixed-size (4KB) key/value block-style storge system in which a key/value pair is stored on a single block. Keys have been limited to a maximum length of 255 bytes; both keys and values are implemented as a []byte. It utilizes memory mapping directly to a backing file for quick reads and writes. It also attempts to re-use deleted blocks and automatically grows the backing file exponentially as it fills.

> 'ads/ohm' is a hashtable implementation that preserves order lexigraphically. It requires you to supply a comparator function (func(i, j) byte) as well as your custom key and value types.
> 'ads/shm' is a thread-safe hashtable that enforces mutual exclusion through sharding and locking of the buckets. The sharding level can be configured in order keep the hashtable quick and also space efficient under different levels of concurrency or paralellism.
> 'ads/xdb' is a document based (nosql) database engine that can be run as a simple embedded solution, or alternativly as a fully database server utilizing the client driver, or binary protocol. It utilizes a basic binary protocol for transmiting data over TCP as well as for storing documents. It has three persistence modes: in memory only (none), timed snapshots (eventual), or transactional (consistent) to choose from. The key has been limited to a maximum length of 255 bytes.
