# ce-racing
Demo repo for CloudEvents with race data.

Flow:

```
dump-to-ce-kafka -> raw-data-processor
```

## dump-to-ce-kafka

This application takes the raw data bump and push CloudEvents to Kafka

Requirements: 

* [Rust setup](https://rustup.rs/)
* Up and running Kakfa cluster to send data
* A topic to send data with 1 partition

To run:

```
cargo run <input_file.dat> <broker_address> <topic_name>
```


## raw-data-processor

This application takes the raw data from Kafka, converts to json and pushes them to TODO

Requirements:

* Golang setup
* Kafka topic where data are contained

To run:

```
go run <broker_address> <input_topic>
```
