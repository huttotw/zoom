# zoom

Zoom is a load testing utility capable of generating different requests according to some template. This should be used to simulate lots of users on a system, or a bunch of random requests at some endpoint.

## Running
To run a load test and send a new KSUID and random number each time:
```
docker run huttotw/zoom:latest zoom -n 10000 -concurrency 10 -url http://localhost:8080 -template '{"id":"{{ ksuid }}","num":{{ intn 100 }}}' -method POST
```

This will send 10,000 `POST` requests, 10 at a time, to `http://localhost:8080` with a payload formatted like:
```json
{"id":"1FDJSGnlAp4C9XTHZQ2XM0u8vFx","num":15}
```

## Flags
| flag            | definition                                                       | required |
|-----------------|------------------------------------------------------------------|----------|
| n               | the number of requests you want to send, omit for infinite       | no       |
| concurrent      | the number of requests to make at one time                       | yes      |
| method          | the HTTP method used to send the request                         | yes      |
| url             | the URL to send the request to                                   | yes      |
| template        | the Go template string that you want to use to form each request | one of template or template-file |
| template-file   | the file that contains your template to form each reqeust        | one of template or template-file |

## Available Functions
| function   | description                                                                   |
|------------|-------------------------------------------------------------------------------|
| `email`    | a random email address                                                        |
| `enum`     | randomly selects one of the arguments given                                   |
| `intn`     | choose a random number between 0 and the given argument                       |
| `ip`       | a random ip address                                                           |
| `ksuid`    | a ksuid from [github.com/segmentio/ksuid](https://github.com/segmentio/ksuid) |
| `string`   | a random string                                                               |
| `url`      | a random url                                                                  |
| `time`     | the current time                                                              |

## Adding new template functions
Since we use `text/template` to form each request, adding new functions to be available at execution time is as simple as adding to the function map in `cmd/zoom/templates.go`. Please see the official documentation for more detail: https://golang.org/pkg/text/template/#FuncMap.

## License
Copyright © 2019 Trevor Hutto

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this work except in compliance with the License. You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
