## Installation

- Create a new GCP project
- Inside your project, create a new Pub/Sub subscription
- Inside the subscription, create a new topic and use the following configuration:
    - Acknowledge deadline > 10 minutes
    - Type: PUSH
    - URL: https://subscriber-dot-<your-project-id>.appspot.com/log-pubsub-message
- In the file `publisher/publisher.go`, change the value of the two constants `PROJECT_ID` and `TOPIC_ID`.
- Deploy the publisher with the command `goapp deploy -application <your-project-id>` -version beta-001 publisher`
- Deploy the subscriber with the command `goapp deploy -application <your-project-id>` -version beta-001 subscriber`
