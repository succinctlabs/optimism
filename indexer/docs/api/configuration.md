## Configuration

Configuration of the indexer is done with a .toml file

- To configure a known chain such as Optimism, Base, or BaseGoerli simply pass in a `preset = CHAIN_ID`.
- For custom chains additional configuration is needed to configure such as l1 and l2 contract addresses

Here is an example indexer.toml

```toml indexer.toml
[chain]
preset = 10

[rpcs]
l1-rpc = "https://eth-goerli.g.alchemy.com/v2/YOUR_API_KEY"
l2-rpc = "https://eth-goerli.g.alchemy.com/v2/YOUR_API_KEY"

[db]
host = "http://localhost"
port = 4321
user = 'postgres'
password = "postgres"

[api]
hostname: "127.0.0.1"
port: 8080

[metrics]
hostname: "127.0.0.1"
port: 7300
```

## Deploying and running indexer

`@eth-optimism/indexer` consists of a single Golang server 

In addition to the app itself, you will also need a Postgres instance
To run or deploy your app a docker container is provided. 

```yaml Example docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:latest
    healthcheck:
      test: [ "CMD-SHELL", "pg_isready -q -U db_username -d db_name" ]
    ports:
      - "5434:5432"
    volumes:
      - postgres_data:/data/postgres

  indexer:
    build:
      context: ..
      dockerfile: indexer/Dockerfile
    healthcheck:
      test: curl localhost:8080/healthz
    environment:
      - 8080:8080
    volumes:
      - /path/to/my/indexer.toml:/indexer/indexer.toml
    depends_on:
      postgres:
        condition: service_healthy
        
volumes:
  postgres_data:

```

