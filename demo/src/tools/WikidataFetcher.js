import axios from 'axios'

export default class WikidataFetcher {
  constructor(lang = 'en', fallbackLang = 'en') {
    this.cache = {}
    this.lang = lang
    this.fallbackLang = fallbackLang
  }

  fetchWikidataItemLabels (ids) {
    ids = Array.isArray(ids) ? ids : [ids]
    const filteredIDs = ids.filter(id => {
      // Filter out empty ids, and ids we already have in cache
      return !!id && !this.cache[id]
    })
    const promises = []

    // Split the fetch IDs array since the API  limit is 50
    this._splitArray(filteredIDs, 40)
      .forEach(idArray => {
        promises.push(this._getWikidataQueryForIDs(idArray))
      })

    return Promise.all(promises)
      .then(values => {
        const result = {}
        // Merge all values in to one big object
        let allValues = {}
        values.forEach(v => {
          if (!allValues.error) {
            // skip errors from the API for now
            allValues = { ...allValues, ...v }
          }
        })

        // Clean up (transform from complex object to key->value of qID=>string)
        // since it's all already in the requested language anyways
        const cleanValues = {}
        Object.keys(allValues).forEach(qID => {
          cleanValues[qID] = (
            // requested lang
            (allValues[qID].labels[this.lang] && allValues[qID].labels[this.lang].value) ||
            // fallback lang
            (allValues[qID].labels[this.fallbackLang] && allValues[qID].labels[this.fallbackLang].value) ||
            // qid
            qID
          )
        })

        // Store in cache
        this.cache = { ...this.cache, ...cleanValues }

        // Merge back with requested IDs from cache (that we removed above)
        ids.forEach(requestedID => {
          if (!requestedID) {
            // skip empty values, just in case
            return
          }
          result[requestedID] = this.cache[requestedID]
        })

        return result
      })
  }

  // private
  _getWikidataQueryForIDs(ids) {
    return axios
      .get(
        'https://www.wikidata.org/w/api.php',
        {
          params: {
            action: 'wbgetentities',
            props: 'labels',
            ids: ids.join('|'),
            languages: this.lang + '|' + this.fallbackLang,
            origin: '*',
            format: 'json'
          }
        },
        {
          headers: {
            'Content-Type': 'application/json; charset=UTF-8',
            origin: '*'
          }
        }
      )
      .then(result => {
        if (result.data.error) {
          return Promise.reject(result.data.error)
        }

        return result.data.entities
      })
      .catch(error => {
        console.log('error', error)
        return {}
      })
  }

  _splitArray(arr, chunkSize = 10) {
    let i, j
    const result = []
    for (i = 0, j = arr.length; i < j; i += chunkSize) {
      result.push(arr.slice(i, i + chunkSize))
    }

    return result
  }
}
