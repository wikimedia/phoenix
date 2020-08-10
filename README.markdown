Phoenix
=======

# Terminology
Define the terms you'll use in this statement

# What?
What is the problem or opportunity? What does it look like if we made this decision and executed on it?

# Why?
Why is this valuable? What organizational objective does this support? And how? What if we do nothing?

# Who?
Can you calculate the RACI? Who is working on the project directly? Who is maintaining a product or code which will be impacted? Who has an operational responsibility which will be impacted?

## Responsible
- Accountable (Who is overseeing the work)
- Consulted (Who do you need to talk to and why. How do they intersect with this problem or - opportunity? Why consulted? What is the level of impact?)
- Informed (For example, program managers/project with dependencies on this work)

# When
What is the timeframe for making this decision. Are there already known milestones?

# Models
Visual representation of the problem/opportunity

# Unresolved Questions
Specific questions that would need to be answered during this decision process

# Resources
[Initial Executive Summary](https://docs.google.com/document/d/1lS9V_knDSIA2Boyax93BFW6MmPlfTAQ-TdI-QshoOHU)
[Model your application domain, not your JSON structures](http://www.markus-lanthaler.com/research/model-your-application-domain-not-your-json-structures.pdf)

# Content descriptions
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
  </tbody>
</table>
