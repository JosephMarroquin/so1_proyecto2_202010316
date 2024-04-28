# Joseph Jeferson Marroquin Monroy 202010316

**Introducción:**
Este proyecto se enfoca en la implementación de un sistema de votaciones para un concurso de bandas de música guatemalteca. Se emplean diversas tecnologías y servicios, como Kubernetes, Kafka, Redis, MongoDB, Grafana y Cloud Run, para gestionar el flujo de datos y proporcionar análisis en tiempo real.

**Objetivos:**
El objetivo principal es desarrollar un sistema eficiente y escalable que permita registrar y analizar las votaciones del concurso de bandas. Esto implica utilizar tecnologías específicas para el procesamiento de datos, almacenamiento y visualización.

**Descripción de cada tecnología utilizada:**
1. **Locust:** Generador de tráfico desarrollado en Python para enviar datos a los servidores desplegados en Kubernetes.
2. **Kubernetes:** Se utiliza un clúster en Google Cloud para desplegar productores, consumidores y el servidor Kafka. Los servicios se organizan en namespaces.
3. **gRPC:** Se emplea para comunicación entre productores y Kafka, implementado en Golang.
5. **Kafka:** Servidor para encolar datos recibidos de los productores antes de ser procesados por el consumidor.
6. **Redis:** Base de datos para almacenar contadores en tiempo real de los consumidores.
7. **MongoDB:** Base de datos para almacenar logs consultados mediante una aplicación web.
8. **Grafana:** Sistema de dashboards conectado a Redis para visualizar las votaciones en tiempo real.
9. **Cloud Run:** Plataforma para desplegar una API en Node y una webapp en Vue JS para visualizar los registros de MongoDB.

**Descripción de cada deployment y service de K8S:**
- Se despliegan pods para  gRPC, Wasm, Kafka, consumidor, Redis y MongoDB, organizados en distintos namespaces según la función.
- Se configuran servicios para permitir la comunicación entre los distintos componentes.


**Conclusiones:**
- Concluimos que el servicio gRPC mostró un rendimiento superior en términos de tiempo de respuesta en comparación con el servicio Rust en este proyecto. Esta conclusión se basa en la observación de los tiempos de procesamiento de los datos enviados por ambos servicios a Kafka. La eficiencia y la baja latencia de gRPC son atribuibles a su arquitectura basada en protocol buffers y su implementación optimizada en Golang. Esto lo hace especialmente adecuado para escenarios donde se requiere una comunicación rápida y eficiente entre los distintos componentes de un sistema distribuido, como en este caso de procesamiento de votaciones en tiempo real.
- Se utilizaría gRPC cuando se necesite una comunicación eficiente entre servicios con alta concurrencia y baja latencia. Wasm se emplearía cuando se requiera la portabilidad del código y un alto rendimiento en entornos distribuidos.
