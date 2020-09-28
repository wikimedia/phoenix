<template>
  <div class="fetch-page">
    <h1 class="text-center">Fetch a part of a page</h1>
    <FormFetchPageSctions />
  </div>
</template>

<script>
// @ is an alias to /src
import FormFetchPageSctions from '@/components/FormFetchPageSctions'

export default {
  name: 'FetchPart',
  components: {
    FormFetchPageSctions
  }
}
</script>
