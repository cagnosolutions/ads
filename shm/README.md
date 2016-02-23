## 'ads/shm' is a thread-safe hashtable.

It enforces mutual exclusion through sharding and locking of the buckets. The sharding level can be configured in order keep the hashtable quick and also space efficient under different levels of concurrency or paralellism.
