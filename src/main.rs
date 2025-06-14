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


mod args;
mod logging;

use jni::{objects::{JObject, JValue}, JavaVM};

use crate::args::ArgsProcessor;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cnf_msg = "Error: Main class not found.";

    let cmd_args = std::env::args().collect::<Vec<_>>();    

    let args_proc = ArgsProcessor::with_args(cmd_args);
    let vm = JavaVM::new(args_proc.build()?)?;

    let mut vm_env = vm.attach_current_thread()?;

    if !args_proc.dry_run() {
        let main_class = args_proc.main_class().as_ref().expect(cnf_msg);
        let mcls = vm_env.find_class(main_class).expect(cnf_msg);

        let mid = vm_env.get_static_method_id(&mcls, "main", "([Ljava/lang/String;)V")
            .expect(format!("static method '{}.main' with signature '([Ljava/lang/String;)V' not found.", main_class).as_str());

        let scls = vm_env.find_class("java/lang/String")?;
        let array = vm_env
            .new_object_array(args_proc.app_args().len() as i32, scls, JObject::null())?;

        for (i, arg) in args_proc.app_args().iter().enumerate() {
            let jstr = vm_env.new_string(arg)?;
            vm_env.set_object_array_element(&array, i as i32, jstr)?;
        }

        let arg_array = [JValue::Object(&array).as_jni()];
        unsafe {
            vm_env.call_static_method_unchecked(&mcls, &mid, jni::signature::ReturnType::Primitive(jni::signature::Primitive::Void), &arg_array)?;
        }
    }

    Ok(())
}
