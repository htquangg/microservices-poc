<h1 align=center>Super Mall: A Microserice Demo Application</h1>

This application is the user-facing part of an online shop 
where user can browse items, add them to the cart, and purchase them.
It is intended to aid the demonstration and testing a microservice
and cloud native technologies.

## Technical Stack
- Backend building blocks
  - [go-kit/kit](https://github.com/go-kit/kit)
  - [segmentio/kafka-go](https://github.com/segmentio/kafka-go)
  - [xorm.io/xorm](https://gitea.com/xorm/xorm)
  - [pressly/goose](https://github.com/pressly/goose)
  - Utils
    - [htquangg/di](https://github.com/htquangg/di)
    - [uber-go/zap](https://github.com/uber-go/zap)
    - [stretchr/testify](https://github.com/stretchr/testify)
    - github.com/sony/sonyflake
    - google.golang.org/genproto
    - google.golang.org/grpc
    - google.golang.org/protobuf
- Infrastructure
  - MySQL, Kafka
  - Hashicorp Consul
  - docker and docker-compose

## Super Mall - Orchestration Saga

## Services

| No. | Service              | Web Server                                       | RPC Server                         |
| --- | -------------------- | ------------------------------------------------ | ---------------------------------- |
| 1   | customer service     | [http://127.0.0.1:30001](http://127.0.0.1:30001) | [127.0.0.1:30002](127.0.0.1:30002) |
| 2   | notification service | [none](none)                                     | [127.0.0.1:31002](127.0.0.1:31002) |
| 3   | store service        | [http://127.0.0.1:32001](http://127.0.0.1:32001) | [127.0.0.1:32002](127.0.0.1:32002) |
| 4   | order service        | [http://127.0.0.1:33001](http://127.0.0.1:33001) | [127.0.0.1:33002](127.0.0.1:33002) |

## Starting project

## Roadmap

- ✅ Enhance project structure with DDD patterns
- Add testing
- Add and integrate with observability libs and tools
- Add user identity management (authentication and authorization)
- Add resiliency

## Credits

- [microservices-demo](https://github.com/GoogleCloudPlatform/microservices-demo)
- [sock-shop](https://github.com/microservices-demo/microservices-demo)
- [Intelli-Mall](https://github.com/LordMoMA/Intelli-Mall)
- [go-micro/demo](https://github.com/go-micro/demo)
- [wild-workouts-go-ddd-example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [thangchung/go-coffeeshop](https://github.com/thangchung/go-coffeeshop)
