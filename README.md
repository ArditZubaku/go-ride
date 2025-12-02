# GoRide - Microservices

An Uber-style ride-sharing platform built with Go microservices architecture. This project demonstrates modern distributed systems patterns including event-driven architecture, real-time communication, and cloud-native deployment.

## Features

### For Riders
- **Interactive Map Interface**: Click-to-select destination with real-time route visualization
- **Trip Preview**: View route, distance, duration, and pricing before booking
- **Multiple Ride Options**: Choose from Sedan, SUV, Van, or Luxury packages with dynamic pricing
- **Real-time Driver Matching**: Live driver availability and assignment
- **Payment Integration**: Secure payment processing via Stripe
- **Trip Status Updates**: Real-time updates via WebSocket connections

### For Drivers
- **Driver Registration**: Register with vehicle type and location
- **Real-time Location Tracking**: Share location updates using geohash-based positioning
- **Trip Requests**: Receive and respond to trip requests in real-time
- **Package Selection**: Choose which ride types (Sedan, SUV, Van, Luxury) to accept

### Technical Features
- **Event-Driven Architecture**: Asynchronous communication via RabbitMQ
- **gRPC Services**: High-performance inter-service communication
- **WebSocket Support**: Real-time bidirectional communication
- **Route Calculation**: Integration with OSRM for accurate routing
- **Geohash-based Proximity**: Efficient driver-rider matching
- **Kubernetes Deployment**: Containerized services with auto-scaling capabilities

## Architecture

### System Overview

The application follows a microservices architecture with the following components:

```
┌─────────────┐
│   Web UI    │ (Next.js + React + TypeScript)
│  (Frontend) │
└──────┬──────┘
       │ HTTP/WebSocket
┌──────▼──────────────────────────────────────┐
│         API Gateway                          │
│  (HTTP Server + WebSocket Handler)           │
└───┬──────────────┬──────────────────────────┘
    │ gRPC         │ gRPC
┌───▼──────┐  ┌───▼──────────┐
│   Trip   │  │   Driver     │
│ Service  │  │   Service    │
└───┬──────┘  └───┬──────────┘
    │             │
    └──────┬──────┘
           │ RabbitMQ
    ┌──────▼──────┐
    │  RabbitMQ   │
    │  (Events)   │
    └─────────────┘
```

### Services

#### 1. **API Gateway** (`services/api-gateway/`)
- **Port**: `8081`
- **Responsibilities**:
  - HTTP REST API endpoints
  - WebSocket connections for real-time updates
  - Request routing to backend services
  - CORS handling
- **Endpoints**:
  - `POST /trip/preview` - Preview trip route and pricing
  - `POST /trip/start` - Create a new trip
  - `WS /ws/riders` - WebSocket for rider updates
  - `WS /ws/drivers` - WebSocket for driver updates

#### 2. **Trip Service** (`services/trip-service/`)
- **Port**: `9083` (gRPC)
- **Responsibilities**:
  - Trip creation and management
  - Route calculation via OSRM API
  - Fare calculation for different vehicle types
  - Trip status management
  - Event publishing
- **Architecture**: Follows Clean Architecture principles (Domain, Service, Infrastructure layers)

#### 3. **Driver Service** (`services/driver-service/`)
- **Port**: `9082` (gRPC)
- **Responsibilities**:
  - Driver registration/unregistration
  - Driver location tracking
  - Geohash-based driver matching
  - Trip request distribution

#### 4. **Web Frontend** (`web/`)
- **Tech Stack**: Next.js 15, React 19, TypeScript, Tailwind CSS
- **Features**:
  - Interactive Leaflet maps
  - Real-time WebSocket connections
  - Stripe payment integration
  - Responsive design

### Communication Patterns

- **Synchronous**: gRPC for service-to-service communication
- **Asynchronous**: RabbitMQ for event-driven messaging
- **Real-time**: WebSockets for client-server bidirectional communication

## Trip Creation Flow

[![](https://mermaid.ink/img/pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ?type=png)](https://mermaid.live/edit#pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ)

For detailed architecture documentation, see:
- [Trip Creation Flow v1](./docs/architecture/trip-creation-flow-v1.md)
- [RabbitMQ Flow v1](./docs/architecture/rabbitmq-flow-v1.md)

## Technology Stack

### Backend
- **Go 1.24.4** - Core language
- **gRPC** - Inter-service communication
- **Protocol Buffers** - Service contracts
- **RabbitMQ** - Message broker for event-driven architecture
- **MongoDB Driver** - Database abstraction (currently using in-memory storage)

### Frontend
- **Next.js 15** - React framework
- **React 19** - UI library
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Leaflet** - Interactive maps
- **Stripe** - Payment processing

### Infrastructure
- **Docker** - Containerization
- **Kubernetes** - Orchestration
- **Tilt** - Development environment
- **OSRM** - Open Source Routing Machine for route calculation

### Development Tools
- **Protobuf** - Protocol buffer compiler
- **Make** - Build automation

## Project Structure

```
ride-sharing-golang-microservices/
├── services/
│   ├── api-gateway/          # HTTP gateway and WebSocket handler
│   ├── driver-service/      # Driver management service
│   └── trip-service/        # Trip management service
├── web/                      # Next.js frontend application
├── shared/                   # Shared libraries and contracts
│   ├── contracts/           # Shared message contracts
│   ├── messaging/           # RabbitMQ utilities
│   ├── proto/               # Generated protobuf code
│   └── types/               # Shared types
├── proto/                    # Protocol buffer definitions
├── infra/                    # Infrastructure as code
│   ├── development/        # Dev environment configs
│   └── production/          # Production configs
├── docs/                     # Architecture documentation
└── tools/                    # Utility scripts
```

## Getting Started

### Prerequisites

- **Docker** - For containerization
- **Go 1.24.4+** - For building services
- **Tilt** - For local development
- **Kubernetes Cluster** - Minikube, Docker Desktop, or cloud provider
- **Node.js 20+** - For frontend development
- **protoc** - Protocol Buffer compiler

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd ride-sharing-golang-microservices
   ```

2. **Generate Protocol Buffers**
   ```bash
   make generate-proto
   ```

3. **Install frontend dependencies**
   ```bash
   cd web
   npm install
   cd ..
   ```

4. **Start Kubernetes cluster** (if using Minikube)
   ```bash
   minikube start
   ```

### Running the Application

#### Development Mode (Recommended)

Start all services using Tilt:

```bash
tilt up
```

Tilt will:
- Build Docker images for all services
- Deploy to Kubernetes
- Set up port forwarding
- Watch for changes and auto-rebuild

#### Manual Deployment

1. **Build services**
   ```bash
   # Build individual services
   cd services/api-gateway && go build -o ../../build/api-gateway
   cd services/driver-service && go build -o ../../build/driver-service
   cd services/trip-service && go build -o ../../build/trip-service
   ```

2. **Deploy to Kubernetes**
   ```bash
   kubectl apply -f infra/development/k8s/
   ```

3. **Run frontend locally**
   ```bash
   cd web
   npm run dev
   ```

### Monitoring

#### Check Pod Status
```bash
kubectl get pods
```

#### View Logs
```bash
# API Gateway
kubectl logs -f deployment/api-gateway

# Trip Service
kubectl logs -f deployment/trip-service

# Driver Service
kubectl logs -f deployment/driver-service
```

#### Kubernetes Dashboard
```bash
minikube dashboard
```

#### Port Forwarding (if not using Tilt)
```bash
# API Gateway
kubectl port-forward service/api-gateway 8081:8081

# RabbitMQ Management UI
kubectl port-forward service/rabbitmq 15672:15672
```

## Development

### Generating Protocol Buffers

After modifying `.proto` files:

```bash
make generate-proto
```

### Environment Variables

Services use environment variables for configuration. Key variables:

- `HTTP_ADDR` - API Gateway HTTP address (default: `:8081`)
- `RABBITMQ_URI` - RabbitMQ connection string
- `RABBITMQ_DEFAULT_USER` - RabbitMQ username
- `RABBITMQ_DEFAULT_PASS` - RabbitMQ password

See `shared/env/` for all environment variable definitions.

### API Endpoints

#### REST API (API Gateway)

- `POST /trip/preview` - Preview trip with route and pricing
  ```json
  {
    "userID": "string",
    "startLocation": { "latitude": 0.0, "longitude": 0.0 },
    "endLocation": { "latitude": 0.0, "longitude": 0.0 }
  }
  ```

- `POST /trip/start` - Create a new trip
  ```json
  {
    "rideFareID": "string",
    "userID": "string"
  }
  ```

#### WebSocket Endpoints

- `WS /ws/riders?userID=<id>` - Rider real-time updates
- `WS /ws/drivers?userID=<id>&packageSlug=<slug>` - Driver real-time updates

### Ride Packages

The system supports four vehicle types:

| Package | Base Price (cents) | Description |
|---------|-------------------|-------------|
| SUV     | 200               | Spacious ride for groups |
| Sedan   | 350               | Economic and comfortable |
| Van     | 400               | Perfect for larger groups |
| Luxury  | 1000              | Premium experience |

Prices are calculated dynamically based on route distance.