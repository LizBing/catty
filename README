# Catty, A JVM Launcher written in Rust  

## In memory of one of my best friends Catherine Young.

I am writing a JVM in rust recently. I found that the JVM Launcher of OpenJDK
is specified for HotSpot VM, so I would like to implement one that restrictly
supports the standard.  

## Building
Just run `cargo build --release` directly!

The environment variable `PRINT_JAVA_VERSION` is set to construct the content
of argument `-version`. For example: 
```bash
export PRINT_JAVA_VERSION="'21ga'"
cargo build

cargo run -- -version

# It outputs
java version '21ga'
Catty version "0.1.0"
Catty, a Universal Java Virtual Machine launcher, designed by Lei Zaakjyu.
In memory one of my best friends Catherine Young.
```

If not:
```bash
java version unknown
Catty version "0.1.0"
Catty, a Universal Java Virtual Machine launcher, designed by Lei Zaakjyu.
In memory one of my best friends Catherine Young.
```

Have a Nice Day!
