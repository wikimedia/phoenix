Phoenix
=======
## Goal
Explore emerging patterns for building systems that manage and distribute content<sup>[[1]](#ref1)</sup>.

The jump from wiki software to modern web architecture (for example, a [reactive](https://www.reactivemanifesto.org/) system) is huge and difficult to conceptualize. Experimenting with core patterns enables us to gain insight into the possibilities and potential impossibilities. To learn more about the architecture process, read the [artifact](https://www.mediawiki.org/wiki/Architecture_Repository/Strategy/Goals_and_initiatives/Structured_content_proof_of_value).

## Primary focus
**Canonical data modeling**<sup>[[2]](#ref2)</sup><sup>[[3]](#ref3)</sup> allows content to be understood by people, programs and machines outside the traditional boundaries of MediaWiki. And, as far as possible, allows consumers<sup>[[4]](#ref4)</sup> to request only what they need.

Building this model requires defining boundaries around parts and their interrelationships. For example, a page has parts (sections) and is also part of collections (about the same topic). Our work here includes:

- Define a predictable structure<sup>[[5]](#ref5)</sup> using industry-standard formats like schema.org (to support predictability and reusability)
- Break down prexisting structures (all the content on the Philadephia page) into parts (a section on the History of Philadephia) and establish interrelationships between the parts (to support "only what they need") using hypermedia linking.<sup>[[6]](#ref6)</sup>
- Enhance the structure with contextual information by associating parts with Wikidata (to enable natural collections like US Cities) and indexing collections with [Elasticsearch](https://www.elastic.co/elasticsearch).
- Enable interaction with the structure via [API calls](https://graphql.org/). Multiple API calls can be wrapped into a single payload -- or not.

[Working draft of our CDM](https://docs.google.com/spreadsheets/d/1ZWuczQQ0XpzCYS92PKXpIP3FM4ds0XPQyz7q9xR5GuE)

*Note: Honestly, we don't know if it's humanly possible to "structure" Wikipedia documents. We are identifying the Biggest Challenges so we can raise them and resolve them organizationally.* 

## Secondary focus

**Loose coupling:**<sup>[[7]](#ref7)</sup> New ways to interact with, enhance or process content (capabilities) operate independently and are built on top of (or adjacent to) the data model.

**Event-based interactions:**<sup>[[8]](#ref8)</sup> activities in the system happen only when they need to happen (asynchronously) with only the information they need to accomplish their aim.

**CQRS**:<sup>[[9]](#ref9)</sup> The current structure inside of MediaWiki is left alone. When changes happen in MW, the new system reacts by getting the necessary information and translating it into the canonical data model.

[Initial Model of how this works](https://app.lucidchart.com/documents/view/f283e649-cdb6-4275-9452-7114571a82e7/Q3nNnx6PpfFM)

## Implementation
We have purposefully defined the implementation toolset for this PoV. We did this so we can focus on the patterns, which present signification challenges, *before deciding which are valuable long term.* By working in AWS, we don't need to build that toolset.

Our next step is not to put this exact toolset into production. Our next step is to collectively design an infrastructure that supports the value while also considering the tradeoffs.

## What about editing?
We are beginning to model editing events now. This approach is designed to work *with* any CMS, as part of a system, not replace it. We will be carefully thinking through which patterns apply, which don't, and which need to be included that aren't here.

## Upcoming use case
We've modeled a number of use cases, some in partnership with the Structured Data Across Wikipedia team, that will benefit from these patterns. When our initial work is complete, we will be prototyping with the Mobile team, specifically focused on References.

## Definitions and resources
- <a name="ref1">[1]</a>: **Content** is the free knowledge / information being shared on a wiki page about a subject.[↩](#goal)
- <a name="ref2">[2]</a>: [Enterprise Integration Patterns: CDM](https://www.enterpriseintegrationpatterns.com/patterns/messaging/CanonicalDataModel.html)[↩](#primary-focus)
- <a name="ref3">[3]</a>: [CDM pitfalls and approaches](https://www.innoq.com/en/blog/thoughts-on-a-canonical-data-model/)[↩](#primary-focus))
- <a name="ref4">[4]</a>: **Consumers** are people, programs (like a front-end application) and machines who are using content outside the context of MediaWiki. MediaWiki could also be a consumer.[↩](#primary-focus)
- <a name="ref5">[5]</a>: [Model your Application Domain, not your JSON structures](https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.1066.5369&rep=rep1&type=pdf)[↩](#primary-focus)
- <a name="ref6">[6]</a>: [Hypermedia linking](https://www.narwhl.com/hypermedia-linking)[↩](#primary-focus)
- <a name="ref7">[7]</a>: **Loose coupling** is an approach to interconnecting the components in a system so that those componentsdepend on each other to the least extent practicable.[↩](#secondary-focus)
- <a name="ref8">[8]</a>: [Event-driven architecture on English Wikipedia](https://en.wikipedia.org/wiki/Event-driven_architecture)[↩](#secondary-focus)
- <a name="ref9">[9]</a>: Command Query Responsibility Segregation (**CQRS**) means that the data model for reading doesn't have to be the same as the model for updating.[↩](#secondary-focus)

## Content descriptions
<table>
  <thead>
    <tr>
      <th></th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td nowrap><code>common</code></td>
      <td>Common structures, helpers, etc</td>
    </tr>
    <tr>
      <td nowrap><code>env</code></td>
      <td>Package for project-wide constants (AWS account &amp; resource information)</td>
    </tr>
    <tr>
      <td nowrap><code>event-bridge</code></td>
      <td>Send filtered change events to an SNS topic</td>
    </tr>
    <tr>
      <td nowrap><code>lambdas/fetch-changed</code></td>
      <td>Subscribe to change events and download the corresponding Parsoid HTML to an S3</td>
    </tr>
    <tr>
      <td nowrap><code>lambdas/fetch-schema.org</code></td>
      <td>Create schema.org JSON-LD output from Wikidata, and upload to S3. Triggered when HTML is added to <code>incoming/</code> (see <code>lambdas/fetch-changed</code>)</td>
    </tr>
    <tr>
      <td nowrap><code>lambdas/merge-schema.org</code></td>
      <td>Merge JSON-LD with HTML documents, and upload to S3. Triggered when linked data is added to <code>schema.org/</code> (see <code>lambdas/fetch-schema.org</code>)</td>
    </tr>
    <tr>
      <td nowrap><code>iac</code></td>
      <td>Terraform configuration for deploying all lambdas to AWS</td>
    </tr>
  </tbody>
</table>

## Terraform IaC

The `iac` directory contain a configuration for deploying lambda function to AWS.
- First setup `~/.aws/credentials` and `~/.aws/config`.
  This should be an [AWS API user with administrative access](https://docs.aws.amazon.com/IAM/latest/UserGuide/getting-started_create-admin-group.html).
- Then provide a proper configuration in `iac/variables.tf`:
```hcl
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "YOUR_REGION_NAME"
}

variable "aws_profile" {
  description = "AWS profile name"
  type        = string
  default     = "YOUR_PROFILE_NAME"
}
```

- Configure an s3 bucket and DynamoDB table for storing TF state and locking it.
  (See [Terraform s3 backend docs](https://www.terraform.io/docs/language/settings/backends/s3.html))
- Provide proper values to TF backend configuration in `iac/main.tf`:

```hcl
terraform {
  ...
  backend "s3" {
    bucket = "STATE_BUCKET_NAME"
    key    = "iac/state"
    region = "YOUR_REGION_NAME"
    dynamodb_table = "LOCK_TABLE_NAME"
    profile = "YOUR_PROFILE_NAME"
  }

```
For more details of the usage see `iac/shared/lambda/README.md`.

*IMPORTANT NOTE*: This IaC DOESN'T create any s3 buckets, SNS Topics and DynamoDB tables. 
Please make sure all required resources have been created in the AWS management console.

[Terraform CLI documentation](https://www.terraform.io/docs/cli/index.html)
