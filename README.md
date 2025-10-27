# solace-logstash-opensearch

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Setup Guide](#setup-guide)
  - [1. Clone the Repository](#1-clone-the-repository)
  - [2. Solace PubSub+ Event Broker Configuration](#2-solace-pubsub-event-broker-configuration)
  - [3. Logstash Configuration](#3-logstash-configuration)
  - [4. Run with Docker Compose](#4-run-with-docker-compose)
- [Configuration Details](#configuration-details)
  - [Solace Connection](#solace-connection)
  - [Logstash Pipeline (`logstash/pipeline/logstash.conf`)](#logstash-pipeline-logstashpipelinelogstashconf)
  - [OpenSearch Connection](#opensearch-connection)
- [Usage](#usage)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Overview

This project provides a robust solution for ingesting real-time event streams from Solace PubSub+ Event Broker into OpenSearch for powerful analytics, search, and visualization. It leverages Logstash as an intermediary data processing pipeline, allowing for flexible data transformation and enrichment before indexing into OpenSearch.

## Features

- **Real-time Ingestion**: Capture Solace messages as they happen.
- **Flexible Data Transformation**: Utilize Logstash's rich set of plugins for parsing, filtering, and enriching data.
- **Centralized Logging & Analytics**: Store and analyze your Solace event data in OpenSearch.
- **Visualization**: Create dynamic dashboards and visualizations using OpenSearch Dashboards.
- **Containerized Deployment**: Easy setup and management using Docker and Docker Compose.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: Usually comes with Docker Desktop. If not, [install Docker Compose](https://docs.docker.com/compose/install/).
- **A running Solace PubSub+ Event Broker**: This can be a local Docker instance, a cloud service, or a hardware appliance. This guide assumes you have access to one.

## Setup Guide

### 1. Clone the Repository

First, clone this repository to your local machine:

```bash
git clone https://github.com/your-org/solace-logstash-opensearch.git
cd solace-logstash-opensearch
```

```bash
docker-compose up -d --build
```

## Usage

1.  **Verify Services**: Ensure all services are running:

    ```bash
    docker-compose ps
    ```

    You should see `opensearch`, `opensearch-dashboards`, and `logstash` in a healthy state.

2.  **Access OpenSearch Dashboards**:
    Open your web browser and navigate to `http://localhost:5601`.

    - **Username**: `admin`
    - **Password**: `admin` (or whatever you set in `docker-compose.yml`)

3.  **Create an Index Pattern**:

    - In OpenSearch Dashboards, go to **Stack Management** -> **Index Patterns**.
    - Click **Create index pattern**.
    - Enter `solace-logs-*` as the index pattern.
    - Select `@timestamp` as the time field.
    - Click **Create index pattern**.

```bash
curl -X GET -k -u admin:admin "https://localhost:9200/solace-logs-*/_search?pretty"

    curl -X GET -k -u admin:admin "https://localhost:9200/_cat/indices?v"
```

4.  **Explore Data**:

    - Go to **Discover** in OpenSearch Dashboards.
    - You should now see your Solace messages being indexed and visualized.

5.  **Publish Messages to Solace**:
    Use any Solace client (e.g., Solace PubSub+ Manager, `sdkperf`, a custom application) to publish messages to the queue configured for Logstash (e.g., `logstash_queue`). Ensure the messages are in JSON format for optimal parsing by the provided Logstash configuration.

    go run publish/publish.go

    **Example JSON Message:**

    ```json
    {
      "timestamp": "2023-10-27T10:30:00.123Z",
      "level": "INFO",
      "source": "my-application",
      "message": "User 'john.doe' logged in successfully.",
      "transactionId": "abc-123-xyz",
      "data": {
        "userId": "john.doe",
        "ipAddress": "192.168.1.100"
      }
    }
    ```

## Troubleshooting

- **Logstash not connecting to Solace**:
  - Check `SOLACE_HOST`, `SOLACE_VPN`, `SOLACE_USERNAME`, `SOLACE_PASSWORD`, and `SOLACE_QUEUE` in `docker-compose.yml`.
  - Verify network connectivity between the Logstash container and the Solace broker. If Solace is on your host, try `host.docker.internal` (Docker Desktop) or your host's IP.
  - Check Solace broker logs for connection errors or authentication failures.
  - Ensure the Solace client username has the necessary permissions (connect, consume from queue).
- **Logstash not sending to OpenSearch**:
  - Check `OPENSEARCH_HOST`, `OPENSEARCH_PORT`, `OPENSEARCH_USERNAME`, `OPENSEARCH_PASSWORD` in `docker-compose.yml`.
  - Check Logstash container logs (`docker-compose logs logstash`) for errors related to OpenSearch output.
  - Ensure OpenSearch is running and accessible (`http://localhost:9200`).
- **Messages not appearing in OpenSearch Dashboards**:
  - Verify Logstash is running and processing messages (check `docker-compose logs logstash`).
  - Ensure your index pattern `solace-logs-*` is correctly configured in OpenSearch Dashboards.
  - Check the time range in OpenSearch Dashboards Discover tab.
  - Confirm messages are being published to the correct Solace queue.
- **Docker Compose issues**:
  - Ensure Docker and Docker Compose are correctly installed and running.
  - Check for port conflicts if other services are using `5601`, `9200`, `9600`.

## Contributing

Feel free to fork this repository, open issues, and submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
