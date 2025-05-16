# Microservices-Template
the Microservices Template for Golang

## GRPC

### 安装依赖
```bash
go get -u "google.golang.org/grpc" \
        "github.com/grpc-ecosystem/go-grpc-middleware"
```

### 安装Protobuf编译器

下载地址：`https://github.com/protocolbuffers/protobuf/releases`

### 安装 Protobuf-Go 插件
安装 protoc-gen-go 插件, 生成 *.pb.go 文件：包含序列化和反序列化功能
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
安装 protoc-gen-go-grpc 插件：生成 *_grpc.pb.go 文件：包含客户端和服务器端的代码
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 编写proto文件

```proto
syntax = "proto3";
package user;

message User {
    string name = 1;
    int32 age = 2;
}
```

### 生成 *.pb.go 文件
- --go_out=. --go_opt=paths=source_relative 用以在 .proto 文件同目录下生成 goods.pb.go
- --go-grpc_opt=paths=source_relative 用以在 .proto 文件同目录下生成 goods_grpc.pb.go

```bash
# paths 可选：source_relative，import，$PREFIX
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto
```

### GRPC拦截器概念
gRPC 提供了两种类型的拦截器：
1. 一元拦截器 (Unary Interceptors): 用于普通的一元 RPC 调用（一次请求，一次响应）。
  - 服务端: grpc.UnaryServerInterceptor，函数签名类似 func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)。
  - 客户端: grpc.UnaryClientInterceptor，函数签名类似 func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error。
2. 流拦截器 (Stream Interceptors): 用于流式 RPC 调用（客户端流、服务端流、双向流）。
  - 服务端: grpc.StreamServerInterceptor，函数签名类似 func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error。
  - 客户端: grpc.StreamClientInterceptor，函数签名类似 func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)。

## 提供RESTful API服务

### 安装依赖
```bash
go get -u github.com/grpc-ecosystem/grpc-gateway/v2
go get -u github.com/grpc-ecosystem/grpc-gateway/v2/runtime@v2.26.3
```

### 生成样板代码
```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2
protoc -I=proto \
   --go_out=entity --go_opt=paths=source_relative \
   --go-grpc_out=entity --go-grpc_opt=paths=source_relative \
   --grpc-gateway_out=entity --grpc-gateway_opt=paths=source_relative \
   proto/entity.proto
```

## OTEL配置
### 安装依赖
```bash
go get -u "go.opentelemetry.io/otel" \
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" \
  "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric" \
  "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" \
  "go.opentelemetry.io/otel/propagation" \
  "go.opentelemetry.io/otel/sdk/metric" \
  "go.opentelemetry.io/otel/sdk/resource" \
  "go.opentelemetry.io/otel/sdk/trace" \
  "go.opentelemetry.io/otel/semconv/v1.26.0" \
  "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp" \
  "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"


```
### OTEL各个组件作用
OpenTelemetry (OTel) 是一个开源的可观测性框架，用于生成、收集和导出遥测数据（metrics、traces、logs），以帮助你理解软件的性能和行为。它由 API、SDK 和 Exporter 组成。

1.  **`go.opentelemetry.io/otel`**
    *   **作用**: 这是 OpenTelemetry Go 项目的核心 API 包。
    *   **详细解释**:
        *   **API 定义**: 它定义了 OpenTelemetry 的核心接口和数据类型，应用程序代码直接与这些 API 交互来生成遥测数据。这包括：
            *   `TracerProvider`: 用于获取 `Tracer` 实例。
            *   `Tracer`: 用于创建和管理 Spans (追踪的基本单元)。
            *   `MeterProvider`: 用于获取 `Meter` 实例。
            *   `Meter`: 用于创建和管理 Instruments (如 Counters, Gauges, Histograms)。
            *   `propagation.TextMapPropagator`: 用于跨服务边界传播上下文（例如 Trace ID）。
            *   `Baggage`: 用于在请求上下文中携带键值对信息。
            *   `trace.SpanContext`: 包含 Trace ID, Span ID, Trace Flags 等。
        *   **解耦**: 关键在于，这个 API 包本身不包含任何具体的实现。你的应用程序代码依赖于这些稳定的 API 接口，而具体的遥测数据收集、处理和导出逻辑则由 SDK (Software Development Kit) 提供。这种分离允许你更换 SDK 实现或配置而无需修改应用程序的检测代码。
        *   **全局实例**: 它还提供了访问全局 `TracerProvider`、`MeterProvider` 和 `TextMapPropagator` 的方法，方便在整个应用中使用。

2.  **`go.opentelemetry.io/otel/exporters/stdout/stdoutmetric`**
    *   **作用**: 这是一个将 Metrics 数据导出到标准输出 (stdout) 的 Exporter。
    *   **详细解释**:
        *   **Exporter**: 在 OpenTelemetry 中，Exporter 负责将 SDK 处理和聚合后的遥测数据发送到指定的后端或目的地。
        *   **stdoutmetric**: 这个特定的 exporter 将指标数据（如计数器的值、直方图的分布等）以可读的格式（通常是 JSON 或简单的文本）打印到控制台。
        *   **用途**: 主要用于开发、调试或演示目的。它不适合生产环境，因为标准输出不是持久化或可查询的遥测存储。通过它，你可以快速验证你的应用是否正确地生成了指标数据。

3.  **`go.opentelemetry.io/otel/exporters/stdout/stdouttrace`**
    *   **作用**: 这是一个将 Traces 数据导出到标准输出 (stdout) 的 Exporter。
    *   **详细解释**:
        *   **stdouttrace**: 类似于 `stdoutmetric`，这个 exporter 将追踪数据（Spans 的信息，如名称、开始/结束时间、属性、事件等）以可读的格式打印到控制台。
        *   **用途**: 同样主要用于开发、调试或演示。它可以帮助你查看应用生成的 Trace 结构，验证 Span 是否正确关联，以及属性是否按预期记录。

4.  **`go.opentelemetry.io/otel/propagation`**
    *   **作用**: 提供上下文传播 (Context Propagation) 的接口和实现。
    *   **详细解释**:
        *   **分布式追踪的关键**: 在微服务架构中，一个用户请求可能会流经多个服务。为了将这些服务中的操作关联成一个完整的分布式追踪，需要在服务调用之间传递追踪上下文（如 Trace ID, Span ID）。
        *   **Propagators**: `TextMapPropagator` 接口定义了如何将上下文信息注入 (inject) 到传出请求的载体（如 HTTP headers）中，以及如何从传入请求的载体中提取 (extract) 上下文信息。
        *   **标准格式**: 这个包提供了对标准传播格式的支持，如：
            *   `W3CTraceContext`: W3C Trace Context 规范，是推荐的默认格式。
            *   `W3CBaggage`: W3C Baggage 规范，用于传递业务相关的键值对。
            *   `B3`: Jaeger 和 Zipkin 等早期追踪系统使用的格式。
        *   通过配置全局或局部的 Propagator，OpenTelemetry SDK 和自动化仪表库可以自动处理上下文的注入和提取。

5.  **`go.opentelemetry.io/otel/sdk/metric`**
    *   **作用**: 这是 OpenTelemetry Go 项目的 Metrics SDK 实现。
    *   **详细解释**:
        *   **SDK 实现**: 它实现了 `go.opentelemetry.io/otel` 中定义的 `MeterProvider` 和 `Meter` 等接口。
        *   **数据处理**: SDK 负责实际的指标数据收集、聚合、过滤和导出。
        *   **核心组件**:
            *   `MeterProvider`: SDK 的入口点，用于创建 `Meter` 实例，并配置指标的收集和导出管道。
            *   `Reader`: 定义了 SDK 如何以及何时从 instruments 中读取指标数据并将其传递给 Exporter。例如 `PeriodicReader` 会定期拉取数据。
            *   `View`: 允许你自定义指标的聚合方式、名称、描述和属性。
            *   `Exporter`: SDK 使用配置的 Exporter (如 `stdoutmetric` 或生产级的 OTLP exporter) 来发送数据。
        *   **配置**: 你需要实例化并配置这个 SDK，例如设置 `Resource`、`Reader` 和 `Exporter`，然后将其注册为全局 `MeterProvider`。

6.  **`go.opentelemetry.io/otel/sdk/resource`**
    *   **作用**: 用于定义和描述产生遥测数据的实体 (Resource)。
    *   **详细解释**:
        *   **Resource概念**: `Resource` 是一组描述产生遥测数据（traces, metrics, logs）的实体的不变属性。这个实体可以是你的服务实例、一个 Kubernetes Pod、一个主机、一个容器等。
        *   **属性**: 例如，`service.name`, `service.version`, `host.name`, `cloud.provider` 等都是 Resource 属性。这些属性会附加到该 Resource 产生的所有遥测数据上。
        *   **标准化**: 使用标准化的 Resource 属性（定义在 Semantic Conventions 中）有助于在遥测后端进行统一的查询、过滤和聚合。
        *   **创建与合并**: 这个包提供了创建 `Resource` 对象的方法，例如从环境变量、预定义检测器（如主机信息、操作系统信息）或手动指定的属性。它也支持合并多个 Resource。
        *   **用途**: 在初始化 SDK (TracerProvider 或 MeterProvider) 时，通常会配置一个 `Resource` 对象。

7.  **`go.opentelemetry.io/otel/sdk/trace`**
    *   **作用**: 这是 OpenTelemetry Go 项目的 Tracing SDK 实现。
    *   **详细解释**:
        *   **SDK 实现**: 它实现了 `go.opentelemetry.io/otel` 中定义的 `TracerProvider` 和 `Tracer` 等接口。
        *   **数据处理**: SDK 负责实际的 Span 创建、采样决策、处理（添加属性、事件等）和导出。
        *   **核心组件**:
            *   `TracerProvider`: SDK 的入口点，用于创建 `Tracer` 实例，并配置追踪的收集和导出管道。
            *   `SpanProcessor`: 定义了 Span 在完成时如何被处理。例如：
                *   `SimpleSpanProcessor`: Span 完成后立即同步导出。
                *   `BatchSpanProcessor`: 批量异步导出 Spans，性能更好，是生产推荐。
            *   `Sampler`: 决定哪些 Trace 应该被记录和导出（采样）。例如 `AlwaysSample`, `NeverSample`, `TraceIDRatioBased`。
            *   `Exporter`: SDK 使用配置的 Exporter (如 `stdouttrace` 或 OTLP exporter) 来发送数据。
        *   **配置**: 你需要实例化并配置这个 SDK，例如设置 `Resource`、`Sampler`、`SpanProcessor` 和 `Exporter`，然后将其注册为全局 `TracerProvider`。

8.  **`go.opentelemetry.io/otel/semconv/v1.26.0`**
    *   **作用**: 提供 OpenTelemetry Semantic Conventions (语义约定) 的常量。
    *   **详细解释**:
        *   **Semantic Conventions**: 为了确保遥测数据在不同语言、框架和后端之间的一致性和互操作性，OpenTelemetry 定义了一套标准的属性键名和枚举值。这些被称为语义约定。
        *   **常量**: 这个包（`v1.26.0` 表示它遵循 OpenTelemetry Semantic Conventions 规范的 1.26.0 版本）提供了这些标准键名的 Go 常量。例如：
            *   `semconv.ServiceNameKey` (值为 "service.name")
            *   `semconv.HTTPMethodKey` (值为 "http.method")
            *   `semconv.DBStatementKey` (值为 "db.statement")
        *   **好处**: 使用这些常量而不是硬编码字符串可以：
            *   避免拼写错误。
            *   当约定更新时，可以通过更新这个包来保持同步。
            *   提高代码的可读性和可维护性。
        *   **应用**: 在手动创建 Span 或指标时，以及在编写自定义检测时，你应该使用这些常量来设置属性。

9.  **`go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`**
    *   **作用**: 为 Go 标准库 `net/http` 包提供自动仪表化 (Instrumentation)。
    *   **详细解释**:
        *   **Instrumentation**: 指在代码中集成遥测逻辑以捕获数据。它可以是手动的（显式调用 API）或自动的（通过库或代理）。
        *   **`otelhttp`**: 这个包属于 `contrib` (贡献) 仓库，提供了对 Go 的 HTTP 客户端和服务端进行自动追踪的功能。
        *   **客户端**: 它可以包装 `http.Client` (通过 `otelhttp.NewTransport`)，使其发出的每个 HTTP 请求都会自动创建一个 Span，记录请求和响应的详细信息（如 URL, method, status code），并自动注入追踪上下文到请求头中。
        *   **服务端**: 它可以包装 `http.Handler` (通过 `otelhttp.NewHandler`)，为每个接收到的 HTTP 请求自动创建一个 Span，记录相关信息，并从请求头中提取追踪上下文，将新创建的 Span 与上游的 Span 关联起来。
        *   **便捷性**: 使用这类自动化仪表库可以大大减少手动编写追踪代码的工作量，并确保关键组件被正确地追踪。

总结：
*   `otel` 定义了核心 API。
*   `sdk/*` 提供了 API 的具体实现，用于处理和导出遥测数据。
*   `exporters/*` 提供了将数据发送到特定目的地的工具（`stdout` 主要用于开发）。
*   `propagation` 处理分布式追踪中的上下文传递。
*   `semconv` 提供了标准化的属性键名。
*   `contrib/instrumentation/*` 为常用库提供了便捷的自动仪表化。

# 参考链接：
1. https://protobuf.dev/reference/go/go-generated/
2. https://opentelemetry.opendocs.io/docs/instrumentation/go/getting-started
