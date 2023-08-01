# katastasi

is a simple page that reads config maps from a kubernetes cluster and builds out configured status pages from that config

## screenhots
this is how it looks like
### Start Page
![start-page](./screenshots/e2e/start-page.png)

### Sample Status Page
![sample-status-page](./screenshots/e2e/failing-status-page.png)

## configuration
a service can be used in multiple status pages
the service defines the metrics that it wants to be checked on the configured prometheus instances
if a service is not configured in a status page, it will not be checked
additionaly one can add maintanance windows to a service, so that it will not be reported as down during that time
````mermaid
    graph RL
    sx[service x] --> sa[status a]
    sy[service y] --> sb[status b]
    sy --> sa
````

