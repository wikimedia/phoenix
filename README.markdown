Phoenix
=======

<table>
  <thead>
    <tr>
      <th></th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td nowrap><a href="/wikimedia/phoenix/blob/master/common"><code>common</code></a></td>
      <td>Common structures, helpers, etc</td>
    </tr>
    <tr>
      <td nowrap><a href="/wikimedia/phoenix/blob/master/env"><code>env</code></a></td>
      <td>Package for project-wide constants (AWS account &amp; resource information)</td>
    </tr>
    <tr>
      <td nowrap><a href="/wikimedia/phoenix/blob/master/event-bridge"><code>event-bridge</code></a></td>
      <td>Send filtered change events to an SNS topic</td>
    </tr>
    <tr>
      <td nowrap><a href="/wikimedia/phoenix/blob/master/lambdas/fetch-changed"><code>lambdas/fetch-changed</code></a></td>
      <td>Subscribe to change events and download the corresponding Parsoid HTML to an S3</td>
    </tr>
    <tr>
      <td nowrap><a href="/wikimedia/phoenix/blob/master/lambdas/fetch-schema.org"><code>lambdas/fetch-schema.org</code></a></td>
      <td>Create schema.org JSON-LD output from Wikidata, and upload to S3. Triggered when HTML is added to <code>incoming/</code> (see <a href="/wikimedia/phoenix/blob/master/lambdas/fetch-changed"><code>lambdas/fetch-changed</code></a>)</td>
    </tr>
  </tbody>
</table>
