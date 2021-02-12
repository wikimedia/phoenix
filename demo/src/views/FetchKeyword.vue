<template>
  <div class="fetch-page">
    <h1 class="text-center">Fetch sections by keyword</h1>
    <FormFetchKeyword />
  </div>
</template>

<script>
// @ is an alias to /src
import FormFetchKeyword from '@/components/FormFetchKeyword'

export default {
  name: 'FetchPart',
  components: {
    FormFetchKeyword
  }
}
</script>
