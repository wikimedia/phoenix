<template>
  <div class="fetch-page">
    <h1 class="text-center">Fetch page by name</h1>
    <FormFetchByPageName />
  </div>
</template>

<script>
// @ is an alias to /src
import FormFetchByPageName from '@/components/FormFetchByPageName'

export default {
  name: 'FetchPage',
  components: {
    FormFetchByPageName
  }
}
</script>
