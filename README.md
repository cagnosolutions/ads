<div align="center"><img src="https://docs.google.com/drawings/d/1E7Ex2Vc7TpnXffN2nlZ4IEBGSuWwaiAROibe4vuzOZY/pub?w=167&amp;h=114"/><h2 style="color:#db4224;">an ambiguous acronym</h2></div>

#### It stands for 'algorithms and data structures' or 'advanced data structures'</h4>

>It contains data structures, algorithms, and helpers for the go programming language. The items in this package are things that I have come across and needed to use at some time or another that are not contained within the standard library (as of go 1.5.) Most if not all of them have been used in some type of production environment. Most of them should not be considered thread-safe, they will note if they are.

It currently contains the following data structures/algorithms:

- <a href="bpt">'ads/bpt'</a> is a lexigraphically ordered b+tree using []byte for both the internal keys as well as the leaf node values.

- <a href="mio">'ads/mio'</a> is a fixed-size (4KB) key/value block-style storge system in which a key/value pair is stored on a single block. Keys have been limited to a maximum length of 255 bytes; both keys and values are implemented as a []byte. It utilizes memory mapping directly to a backing file for quick reads and writes. It also attempts to re-use deleted blocks and automatically grows the backing file exponentially as it fills.

- <a href="ohm">'ads/ohm'</a> is a hashtable implementation that preserves order lexigraphically. It requires you to supply a comparator function (func(i, j) byte) as well as your custom key and value types.

- <a href="shm">'ads/shm'</a> is a thread-safe hashtable that enforces mutual exclusion through sharding and locking of the buckets. The sharding level can be configured in order keep the hashtable quick and also space efficient under different levels of concurrency or paralellism.

- <a href="xdb">'ads/xdb'</a> is a document based (nosql) database engine that can be run as a simple embedded solution, or alternativly as a fully stand-alone database server utilizing the driver. It utilizes a basic binary protocol for both transmiting data over the wire as well as for storing the documents. The keys have been limited to a maximum length of 255 bytes.
