/*
 * Copyright (c) 2025, Lei Zaakjyu. All rights reserved.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

use std::{fs::File, io::{BufRead, BufReader}, process::exit};

use jni::{InitArgs, InitArgsBuilder, JvmError};
use zip::ZipArchive;

use crate::{logging::{print_help, print_x}, show_version};

pub struct ArgsProcessor {
    _main_class: Option<String>,
    _class_path: String,

    _jvm_args: Vec<String>,
    _app_args: Vec<String>,

    _jar_mode: bool,
    _module_mode: bool,
    _cp_mode: bool,

    _dry_run: bool
}

impl ArgsProcessor {
    fn new() -> Self {
        ArgsProcessor {
            _main_class: None,
            _class_path: String::from("."),

            _jvm_args: Vec::new(),
            _app_args: Vec::new(),

            _jar_mode: false,
            _module_mode: false,
            _cp_mode: false,

            _dry_run: false
        }
    } 
}

impl ArgsProcessor {
    pub fn main_class(&self) -> &Option<String> {
        &self._main_class
    }

    pub fn app_args(&self) -> &Vec<String> {
        &self._app_args
    }

    pub fn dry_run(&self) -> bool {
        self._dry_run
    }
}

fn exit_for_bad_arg() {
    print_help();
    exit(1);
}

impl ArgsProcessor {
    pub fn with_args<'a>(args: Vec<String>) -> Self {
        let mut this = ArgsProcessor::new();

        let mut index = 1;
        if args.len() > 1 {
            while index <= args.len() {
                let arg = &args[index];

                match arg.as_str() {
                    "-showversion" => {
                        show_version!(eprintln);
                        exit(0);
                    }

                    "-version" | "--show-version" => {
                        show_version!(println);
                        exit(0);
                    }

                    "-help" | "--help" | "-?" => {
                        print_help();
                        exit(0);
                    }

                    "-X" => {
                        print_x();
                        exit(0);
                    }

                    "-cp" | "-classpath" => {
                        if index + 1 >= args.len() {
                            eprintln!("Error: -cp requires 1 argument.");
                            exit_for_bad_arg();
                        }

                        if this._module_mode {
                            eprintln!("Error: -cp is not supported in --module mode");
                            exit(1);
                        }

                        if !this._jar_mode {
                            this._cp_mode = true;

                            index += 1;
                            this._class_path = args[index].clone();
                        }
                    }

                    "--dry-run" => {
                        this._dry_run = true;
                    }
                    
                    "-jar" => {
                        if index + 1 >= args.len() {
                            eprintln!("Error: -jar requires 1 argument.");
                            exit_for_bad_arg();
                        }

                        this._jar_mode = true;
                        if this._cp_mode {
                            this._cp_mode = false;
                        }
                        if this._module_mode {
                            eprintln!("Error: -jar is not supported in --module mode.");
                            exit(1);
                        }

                        index += 1;
                        let (main_class, class_path) = process_jar(&args[index]);
                        this._main_class = Some(main_class);
                        this._class_path = class_path;
                    }

                    "-m" | "--module" |
                    "--module-path" |
                    "--add_modules" |
                    "--list-modules" |
                    "-d" | "--describe-module" |
                    "--validate-modules" => {
                        eprintln!("Error: Module processing is not supported for now(Catty Version {}).", env!("CARGO_PKG_VERSION"));
                        exit(1);
                    }

                    _ => break
                }

                index += 1;
            }
            debug_assert!(this._cp_mode as i32 + this._jar_mode as i32 + this._module_mode as i32 <= 1, "Multiple modes.");

            this._jvm_args.push(format!("-Djava.class.path={}", this._class_path));

            // Process JVM args.
            for n in &args[index..] {
                if n.chars().nth(0).unwrap() == '-' {
                    this._jvm_args.push(n.clone());
                } else { break; }

                index += 1;
            }

            if index <= args.len() {
                if let None = this._main_class {
                    this._main_class = Some(args[index].clone());
                }
            }

            for n in &args[index + 1..] {
                this._app_args.push(n.clone());
            }
        } else {
            exit_for_bad_arg();
        }

        this
    }

    pub fn build(&self) -> Result<InitArgs, JvmError> {
        let mut builder = InitArgsBuilder::new().version(jni::JNIVersion::V8);

        for n in &self._jvm_args {
            builder = builder.option(n);
        }

        builder.build()
    }
}

// left: main class, right: full class path
fn process_jar(jar_path: &String) -> (String, String) {
    let failed_open_msg = "java: Failed to open JAR file.";

    let file = File::open(jar_path).expect(failed_open_msg);
    let mut archive = ZipArchive::new(file).expect(failed_open_msg);

    let mut mani_file = archive
        .by_name("META-INF/MANIFEST.MF")
        .expect("Error: 'META-INF/MANIFEST.MF' not found in JAR file.");
    let reader = BufReader::new(&mut mani_file);

    let mut main_class = None;
    let mut class_path = jar_path.clone();

    for line in reader.lines() {
        let line = line.expect("Error: Failed to read JAR file.");

        if line.starts_with("Main-Class:") {
            main_class = Some(line.trim_start_matches("Main-Class:").trim().to_string());
        } else if line.starts_with("Class-Path:") {
            let class_path_str = line.trim_start_matches("Class-Path:").trim().to_string().replace(" ", ":");
            class_path += format!(":{}", class_path_str).as_str();
        }
    }

    (main_class.expect("Error: Main Class not found in JAR file."), class_path)
}


