# Suspicious IP Checker

The Suspicious IP Checker is a microservice-based application designed to identify and alert on suspicious IP addresses. It leverages VirusTotal for IP reputation checks and uses Apache Kafka for asynchronous communication between services.

## Features

*   **IP Submission API**: A RESTful API endpoint to submit IP addresses for analysis.
*   **VirusTotal Integration**: Checks submitted IP addresses against VirusTotal for known threats and reputation.
*   **Kafka Messaging**: Utilizes Kafka for reliable and scalable communication between the IP submission service and other potential services (e.g., an alert service).
*   **Microservice Architecture**: Designed with separate services for IP submission and alerting, promoting scalability and maintainability.
*   **Dockerized Deployment**: Easy to deploy and manage using Docker and Docker Compose.

## Architecture

The application consists of the following main components:

*   **`ip-submission-service`**:
    *   Provides a REST API (`/submit-ip`) to receive IP addresses.
    *   Uses the `service` package to check the IP's reputation (e.g., via VirusTotal).
    *   Publishes the IP analysis results to a Kafka topic.
*   **`alert-service`**:
    *   (Planned/Future) Consumes IP analysis results from Kafka.
    *   Can be extended to trigger alerts (e.g., email, Slack, logging) based on suspicious IP activity.
*   **Kafka**: A distributed streaming platform used for inter-service communication.
*   **Zookeeper**: Manages Kafka brokers.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   Docker and Docker Compose
*   A VirusTotal API Key

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-repo/suspicious-ip-checker.git
    cd suspicious-ip-checker
    ```

2.  **Create a `.env` file:**

    Create a `.env` file in the root directory of the project and add your VirusTotal API key:

    ```
    VIRUSTOTAL_API_KEY=your_virustotal_api_key_here
    ```

3.  **Start the services:**

    ```bash
    docker-compose up --build
    ```

    This will build the Docker images and start the Zookeeper, Kafka, and `ip-submission-service` containers.

## Configuration

The application can be configured using environment variables, primarily through the `.env` file.

*   `VIRUSTOTAL_API_KEY`: Your API key for VirusTotal.
*   `KAFKA_BROKER`: The address of the Kafka broker (e.g., `kafka:9092` for Docker Compose).

## API Endpoints

### `POST /submit-ip`

Submits an IP address for analysis.

**Request Body:**

```json
{
    "ip": "192.168.1.1"
}
```

**Response:**

```json
{
    "ip": "192.168.1.1",
    "status": "clean"
}
```

The `status` field will indicate the reputation of the IP address (e.g., "clean", "malicious", "suspicious").
