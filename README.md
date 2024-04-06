
# SpaceTraders API Client

This document provides a basic overview of how to use the SpaceTraders API client in your Go projects.

## Installation

To use the SpaceTraders API client, you first need to install the package in your project. Run the following command:

To install the SpaceTraders API client package in your Go project, execute the following command in your terminal:

```bash
  go get github.com/jjkirkpatrick/spacetraders-client
```

## Usage

To set up a new client for interacting with the SpaceTraders API, follow these steps:

### Creating a New Client Instance

1. **Import the SpaceTraders client package** into your Go file where you intend to use the client.

2. **Initialize the client**: Create a new instance of the SpaceTraders client by providing the base URL of the SpaceTraders API.

```go
options := client.DefaultClientOptions()
client, cerr := client.NewClient(options, nil)
```

### Making Requests


```go
  // All endpoints from the API are represented with a functionm the returns struct match that as the return model in the API documentation
	dock, err := client.DockShip("ship-1")
	if err != nil {
		logger.Fatalf("Failed to dock ship: %v", err)
	}

	logger.Printf("Docked: %+v", dock.Status)
```


### Paginated Requests


```go
  // Create a new paginator for the list of factions endpoint
	paginator, err := client.ListFactions()

  // Fetch all pages of results
	allPages, err := paginator.FetchAllPages()

  // Iterate over the results
	for _, faction := range allPages {
		logger.Printf("Faction: %+v", faction.Symbol)
	}
```

## Metrics Monitoring with InfluxDB and Grafana

To integrate metrics monitoring with your SpaceTraders API client, follow these steps:

### Setting Up the Metrics Client

1. **Initialize the Metrics Client**: First, you need to create an instance of the `MetricsClient` by providing the InfluxDB connection details including URL, token, organization, and bucket name. This client will be responsible for sending metrics data to InfluxDB.

The NewMetricsClient function takes the following parameters: url, token, org, bucket all of type string

```go
	metricsReporter := metrics.NewMetricsClient(
		"http://192.168.1.33:8086",
		"Token",
		"spacetraders",
		"spacetraders",
	)
```
Once you have created the metrics client, you will need to instansiate the spaceTraders client with the metrics client as a parameter

```go
    options := client.DefaultClientOptions()
    client, cerr := client.NewClient(options, metricsReporter)
```

You will now have access to the client.WriteMetric function, as well as the MetricBuilderBuilder

```go
	metric, merr := metrics.NewMetricBuilder().
		Namespace("spacetraders").
		Tag("agent", agent.Symbol).
		Field("credits", agent.Credits).
		Timestamp(time.Now()).
		Build()

	if merr != nil {
		logger.Fatalf("Failed to build metric: %v", err)
	}

	client.WriteMetric(metric)
```

#### HTTP Metrics

By Default the Client automatically sends metrics to grafana for each request it makes to allow for tracking of requests per second and error rates: 

#### Disabling Metrics, 

If you don't wish to use Grafana metrics you may pass nil to the MetricsReporter parameter on the NewClient function

```go
    client, cerr := client.NewClient(options, nil)
```

Internally this will change the interface used to a no-op interface that will not send any metrics to Influx

To set up InfluxDB and Grafana for monitoring and visualizing metrics, you can use the provided `docker-compose` file. This setup allows you to run InfluxDB and Grafana in containers, making it easy to get started without installing each software individually on your system.

### Prerequisites
- Docker and Docker Compose installed on your machine.

### Steps to Setup
1. **Start the Services**: Navigate to the directory containing the `docker-compose.yml` file and run the following command to start InfluxDB and Grafana services:
   ```bash
   docker-compose up -d
   ```
   This command will download the necessary Docker images and start the services in detached mode.

2. **Access Grafana**: Once the services are up, you can access Grafana by opening `http://localhost:3000` in your web browser. The default login credentials are:
   - **Username**: admin
   - **Password**: mysecretpassword

3. **Configure InfluxDB as a DataSource in Grafana**:
   - In the Grafana dashboard, navigate to **Configuration > Data Sources**.
   - Click on **Add data source**, and select **InfluxDB**.
   - Use the following settings to configure the InfluxDB data source:
     - **URL**: http://influxdb:8086
     - **ORG**: The name you used when setting up InfluxDB from the web interface.
     - **Token**: The token you generated when setting up InfluxDB from the web interface.
   - Click **Save & Test** to verify the connection.

4. **Create Dashboards**: Now, you can create dashboards in Grafana to visualize the metrics stored in InfluxDB. Use the Grafana UI to create and customize your dashboards.

### Stopping the Services
To stop the InfluxDB and Grafana services, run the following command in the directory containing your `docker-compose.yml` file:

#### Example docker-compose.yml
```yaml
version: '3'

services:
  influxdb:
    image: influxdb:2.6
    ports:
      - "8086:8086"
    volumes:
      - influxdb_data:/var/lib/influxdb2

  grafana:
    image: grafana/grafana:9.4.7
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
    environment:
      GF_SECURITY_ADMIN_USER: admin
      GF_SECURITY_ADMIN_PASSWORD: mysecretpassword
    depends_on:
      - influxdb

volumes:
  influxdb_data:
  grafana_data:
```
