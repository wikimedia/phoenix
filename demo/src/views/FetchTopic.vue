<template>
  <div class="fetch-topic">
    <h1 class="text-center">Fetch a topic</h1>
    <FormFetchTopic />
  </div>
</template>

<script>
// @ is an alias to /src
import FormFetchTopic from '@/components/FormFetchTopic'

export default {
  name: 'FetchTopic',
  components: {
    FormFetchTopic
  }
}
</script>
