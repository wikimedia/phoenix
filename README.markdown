Phoenix
=======

| | Description |
| ---- | ---- |
| [`common`](common) | Common structures, helpers, etc |
| [`env`](env) | Package for project-wide constants (AWS account & resource information) |
| [`event-bridge`](event-bridge) | Send filtered change events to an SNS topic |
| [`lambdas/fetch-changed`](lambdas/fetch-changed) | Subscribe to change events and download the corresponding Parsoid HTML to an S3 |
| [`lambdas/fetch-schema.org`](lambdas/fetch-schema.org) | Create schema.org JSON-LD output from Wikidata, and upload to S3. Triggered when HTML is added to `incoming/` (see [`lambdas/fetch-changed`](lambdas/fetch-changed)) |
