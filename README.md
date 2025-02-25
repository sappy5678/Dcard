# How to run
At repo root folder, run the following command to start the server

```
docker compose -f Docker-Compose.yml up --build
```
and the server will run on port 8080

# Key Feature
## MultiLayer Storage Architecture
- PostgreSQL ensures data persistence and consistency. 
- Redis provides high-speed caching, significantly improving redirect speed for high traffic.
## Short Code (Short URL ID) Generation Strategy
- Utilizes an auto-increment ID with Base57 encoding to ensure uniqueness and control the length of the short URL.
- Employs a combination of machine ID and auto-increment ID to generate shortcodes ({machineID}-{AutoIncrID}), offering high efficiency and scalability across multiple nodes.
## Handling non-existent shorten URL 
- Using Bloom Filters to filter out non-existent keys, ensures that requests for keys that do not exist do not reach the database.
## Handling access shorten URL simultaneously handling
- Uses Redis to implement a distributed lock, ensuring that even if multiple requests access the same key simultaneously, only one request will interact with the database.

# Trade off
## Short Code (Short URL ID) Generation Strategy
- Using UUIDs
  - Advantages: Simple and convenient.
  - Disadvantages: Requires checking if the UUID has already been generated; if it has, a new one must be created, which is less efficient.
- Using a centralized server
  - Advantages: High efficiency, centralized management, ability to generate shorter URLs, and can adjust length based on actual demand.
  - Disadvantages: Not easily scalable; requires maintaining an additional service, leading to higher costs.
- Machine ID + Auto-increment ID ({machineID}-{AutoIncrID})
  - Advantages: Highest efficiency, easy to scale, simple implementation, generates relatively short URLs that can increase in length as needed.
  - Disadvantages: Requires a centralized service to issue machine IDs, though this service does not have heavy loading, making maintenance costs moderate.
  - Decision: We chose this approach because it offers excellent performance and scalability, generates reasonably short URLs, and the maintenance cost is acceptable.
## Handling non-existent shorten URL 
- Caching Non-Existent Keys
  - Advantages: Low maintenance cost.
  - Disadvantages: With high volumes of garbage requests, Redis might become overwhelmed, leading to eviction of valid cache entries. Additionally, these requests may still reach the database.
- Using Bloom Filters to Filter Non-Existent Keys
  - Advantages: Ensures that requests for non-existent keys do not reach the database.
  - Disadvantages: Higher maintenance cost, requires a dedicated Redis cluster, and necessitates managing backup, restoration, reconstruction, and write checks for the Bloom filter.
  - Decision: We chose this solution because it delivers excellent performance. Although the maintenance cost is higher, it effectively filters out garbage requests.

# Possible Improvement (Not implement in this demo)
## MultiLayer cache
- Cache key within pod
- Cache key in CDN
## Short Code (Short URL ID) Generation Strategy
- Although each pod can generate IDs, a centralized service is still needed to assign machine IDs to each pod. For simplicity, this centralized machine ID service was not implemented.
- A config server is also required to provide the hostname for the short URL service to each pod.
## Handling Non-existent shorten URL 
- While Bloom filters may have a low probability of false positives, the impact is mitigated by the database's own caching mechanisms.
- To thoroughly mitigate this issue, we can cache the keys that resulted in false positives, reducing the chances of hitting the database.
- It's recommended to host the Bloom filter on a separate Redis cluster, not shared with the cache; otherwise, it might get evicted by Redis.
  - We also need to consider the backup and restoration of the Bloom filter, as well as reconstruction if it gets cleared.
  - Since Redis may lose data even after successful writes, it is recommended to perform a secondary check after a few seconds to ensure the key was added to the Bloom filter.


# API Example

## Upload URL API

```bash
curl -X POST -H "Content-Type:application/json" http://localhost:8080/api/v1/urls -d '{ "url": "<original_url>", "expireAt": "2025-02-28T09:20:41Z"}'
```
### Checking
* url is available format
* expireAt is greater than now

### Response

```json
{ "id": "<url_id>", "shortUrl": "http: //localhost:8080/<url_id>" }
```

## Redirect URL API

```bash
curl -L -X GET http://localhost:8080/<url_id> => REDIRECT to original URL
```
### Checking
* url_id is exist, and not expired


# Unit test
```
❯ go test ./... -cover -p=1 -count=1
        github.com/sappy5678/dcard/cmd/api              coverage: 0.0% of statements
ok      github.com/sappy5678/dcard/pkg/domain   1.081s  coverage: 83.3% of statements
        github.com/sappy5678/dcard/pkg/service          coverage: 0.0% of statements
ok      github.com/sappy5678/dcard/pkg/service/shorturl 0.036s  coverage: 0.0% of statements [no tests to run]
ok      github.com/sappy5678/dcard/pkg/service/shorturl/cache   9.078s  coverage: 79.3% of statements
ok      github.com/sappy5678/dcard/pkg/service/shorturl/logservice      0.040s  coverage: 100.0% of statements
ok      github.com/sappy5678/dcard/pkg/service/shorturl/repository      9.495s  coverage: 78.6% of statements
ok      github.com/sappy5678/dcard/pkg/service/shorturl/shortcode       0.330s  coverage: 92.3% of statements
ok      github.com/sappy5678/dcard/pkg/service/shorturl/transport       0.075s  coverage: 91.7% of statements
ok      github.com/sappy5678/dcard/pkg/utl/config       0.046s  coverage: 100.0% of statements
ok      github.com/sappy5678/dcard/pkg/utl/locker       6.479s  coverage: 100.0% of statements
ok      github.com/sappy5678/dcard/pkg/utl/postgres     8.260s  coverage: 83.3% of statements
ok      github.com/sappy5678/dcard/pkg/utl/redis        5.924s  coverage: 83.3% of statements
ok      github.com/sappy5678/dcard/pkg/utl/server       0.045s  coverage: 27.8% of statements
ok      github.com/sappy5678/dcard/pkg/utl/zlog 0.039s  coverage: 100.0% of statements
```