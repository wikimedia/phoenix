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
