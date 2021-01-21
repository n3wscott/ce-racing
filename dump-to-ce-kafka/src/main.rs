use tokio::io::AsyncReadExt;
use tokio::fs::File;
use rdkafka::producer::{FutureProducer, FutureRecord};
use rdkafka::ClientConfig;
use cloudevents::{EventBuilderV10, EventBuilder};
use cloudevents_sdk_rdkafka::{MessageRecord, FutureRecordExt};
use std::time::Duration;
use std::env;

#[tokio::main]
async fn main() {
    let args: Vec<String> = env::args().collect();
    let input_file_location = args.get(1).expect("Input file location");
    if args.len() == 2 {
        // Dump to console
        dump_to_console(input_file_location)
            .await;
        return
    }
    let brokers = args.get(2).expect("Brokers");
    let topic_name = args.get(3).expect("Topic name");

    let producer: &FutureProducer = &ClientConfig::new()
        .set("bootstrap.servers", brokers)
        .set("message.timeout.ms", "5000")
        .create()
        .expect("Producer creation error");

    let mut f = File::open(input_file_location)
        .await
        .expect("Opened file");
    let mut buffer = [0; 328];

    let mut inc: u32 = 0;

    loop {
        let n = f.read(&mut buffer)
            .await
            .expect("Read file");

        if n == 0 {
            println!("EOF, closing");
            return
        }

        let event = EventBuilderV10::new()
            .id(inc.to_string())
            .ty("acdump")
            .source("http://ac-dump-reader.myapp.com/")
            .data("application/octet-stream", buffer.to_vec())
            .build()
            .unwrap();

        let message_record = MessageRecord::from_event(event)
            .expect("error while serializing the event");

        producer
            .send(
                FutureRecord::to(topic_name)
                    .message_record(&message_record)
                    .key(&inc.to_string()),
                Duration::from_secs(10),
            )
            .await
            .expect("Correct production of messages");

        inc += 1;
    }
}

async fn dump_to_console(input_file_location: &str) {
    let mut f = File::open(input_file_location)
        .await
        .expect("Opened file");
    let mut buffer = [0; 328];

    let mut inc: u32 = 0;

    loop {
        let n = f.read(&mut buffer)
            .await
            .expect("Read file");

        if n == 0 {
            println!("EOF, closing");
            return
        }

        let lap_time = u32::from_le_bytes([buffer[40], buffer[41], buffer[42], buffer[43]]);
        let speed = f32::from_le_bytes([buffer[8], buffer[9], buffer[10], buffer[11]]);

        println!("{}: Lap time '{}', Speed '{}'", inc, lap_time, speed);

        inc += 1;
    }
}
