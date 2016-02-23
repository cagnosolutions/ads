# 'ads/mio' is a fixed-size key/value storge system.

> It uses fixed sized blocks (4KB) in which a key/value pair is stored on a single block. Keys have been limited to a maximum length of 255 bytes; both keys and values are implemented as a []byte. It utilizes memory mapping directly to a backing file for quick reads and writes. It also attempts to re-use deleted blocks and automatically grows the backing file exponentially as it fills.
