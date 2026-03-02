# URL Shortener

An exercise inspired by the System Design Interview [book](https://bytebytego.com/courses/system-design-interview/design-a-url-shortener) by Alex Xu.

## Problem

Build/design a high scale URL Shortening service similar to tiny url or bitly

The expected scale is:

- `100_000_000` new urls created per day
- `1_000_000_000` redirect requests per day
- The service will run for 10 years and should be able to store all URLs for that time for a total of 365B rows
- Assuming 200 bytes per row the total storage required would be: `200 bytes * 365B = ~75TB` with some wiggle room

## Features

- Creating a short url
- Redirecting from a short url to a long url
- Expiration policies for links
    - one time link i.e., expires after one visit
    - maximum age link e.g., expires after 30 days
- Telemetry 
    - Counting visits


## Architecture

For this exercise I've opted to use serverless solutions to avoid having to configure a vpc, servers, containers etc. 
and to optimise the cost for myself as serverless solutions usually have a very good free tier rather than optimizing
costs for the actual scale I'm building for. I cover this in more detail in the [cost estimation section](#cost-estimation)

![URL Shortening Service architecture diagram](./docs/system-design.svg)

TODO: Description

### Create url request flow

![Create url request flow diagram](./docs/create-short-url.svg)

TODO: Description

### Resolve shortcode request flow

![Resolve shortcode request flow diagram](./docs/resolve-short-code.svg)

TODO: Description

### Visit aggregation pipeline

![Visit aggregation pipeline diagram](./docs/log-processing-pipeline.svg)

TODO: Description

## Implementation

### Create URL

TODO: Description

Challenge: generating unique short codes 

Attempt 1: Incrementing counter
- sequentiality
- need to rely on an external storage layer in a distributed system

Attempt 2: Generating random number in a 62^n space
- birthday problem

Attempt 3: Snowflake id adjusted for lambda 
- original Snowflake ID is not usable by lambdas
- sequentiality

Final: Snowflake + Feistel Network for removing sequentiality 
- using snowflake to deterministically generate a unique id
- using the feistel network algorithm to scramble the id

### Resolve URL

TODO: Description

Challenge: handling one-time links in a distributed system reliably

Challenge: introducing caching while at the same time having somewhat accurate counts of the visits

### Aggregating visits

TODO: Description

Challenge: Aggregating 1B visits per day 

## Trade offs

## Cost estimation

TODO: Refine

The biggest cost will be:
- CDN
    - The bulk of the price of cloudfront is the 1.1B requests per day that cannot be avoided and this comes out to ~$33'000 per month
    - Additional costs for transfer is another ~$2600
    - Total: ~$36000
    - The greatest cost benefit of using a CDN is that we do not send all of the 1.1B requests to origin, saving on additional cost on other parts of the infrastructure
- API Gateway
    - Assuming a cache-hit rate of 95%, 50M redirect and 100M new url requests will have to be handled by the API Gateway
    - This comes out to ~$4100 
    - If we weren't caching at the CDN it would be ~$30000 + additional lambda and dynamodb charges for handling billions more invocations and
    database queries a month

- Storage
    - Scale
      - Assuming around 200bytes per row * 350B rows over ten years = 70TB of storage required
      - 100M new urls created each day
      - 1B redirects, assuming aggressive caching resulting in 50M reads each day
      - The visits are aggregated into buckets for an additional of 1B Updates in the database for a total of:
      - 1.1B writes and 50M reads so it's a very write heavy application
    
    - On demand dynamodb
       - For us-east-1 dynamodb charges $0.25 per GB for a total of ~$18000/month
       - For 1.1B writes each day we'll have to pay in total $21000/month
       - 50M reads is under $100/month
       - Grand total of ~$39000/month for on-demand Dynamodb

    - provisioned capacity dynamodb
      -  For provisioned capacity in dynamodb:
      -  the storage stays the same
      -  Write capacity for 10K WCU is $15'000 up front for a year and then ~$1000/month
      -  Read capacity for 300 RCU is ~$90 up front for a year and then ~$7/month

    - Some other options would be:
       - RDS that would require an upfront payment for reserved capacity of ~$100000 and then ~$18000/month
       - Aurora which would require an upfront payment for reserved capacity of ~$57000 and then ~$14000/month
       - DocumentDB is ~$17000/month but with a single instance of 32 vCPU's and 256GiB of memory, each additional instance is another ~$3000/month

    - None of these data stores are well suited for this type of data, clickhouse would be a much better choice for analytics at a tenth of the price or even less (OLTP vs. OLAP)

    
[AWS Calculator Estimate](https://calculator.aws/#/estimate?id=28d3f1350fe8a88958a982ee9306ddd125ec1458)
