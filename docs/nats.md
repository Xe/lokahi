# Nats queues

## check.errors

Sent a `xe.github.lokahi.admin.WebhookData` by `webhookworker` when it processes 
a check that failed due to an underlying error, no reply expected.

## check.run

Sent `xe.github.lokahi.Check` and will return `xe.github.lokahi.admin.Health`.
This performs a HTTP request against the URL defined in the Check and uses metadata
about the response to populate the returned Health.

## webhook.egress

Sent `xe.github.lokahi.admin.WebhookData` and no return expected. This sends
webhook data to the URL in the underlying Check.
