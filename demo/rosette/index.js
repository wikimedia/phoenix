import ApiHelper from './ApiHelper.js'
import { createRequire } from 'module';
const require = createRequire(import.meta.url);
const { ArgumentParser } = require('argparse');
const api = new ApiHelper();

var parser = new ArgumentParser({
  add_help: true,
  description: "Fetch data for a section part from rosette"
});

parser.add_argument("--page", {help: "Page name from simple.wikipedia", required: false});
parser.add_argument("--part", {help: "A section name to fetch", required: false});
parser.add_argument("--endpoint", {help: "Endpoint for Rosette API", required: false});
parser.add_argument("--output", {help: "Output type: 'table', 'csv', 'collectconcepts'. Default: table", required: false});
const args = parser.parse_args();
// const rosetteApi = new Api(config.rosette_key, args.url);

// Available endpoints: 'topics', 'entities', 'categories'
const endpoint = args.endpoint || "topics";
const pageName = args.page || 'Philadelphia'
const requestedPart = args.part || ''
const requestedOutput = args.output || 'table'
const queries = {
  parts: `{
    page(name: { authority: "simple.wikipedia.org", name: "${pageName}"} ) {
      name
      dateModified
      hasPart(offset: 0) {
        name
        unsafe
      }
    }
  }`
}
// Fetch from graphql service
console.log('Fetching...')
api.fetchFromGraphQL(pageName).then(parts => {
    let partIndex = Math.floor(Math.random() * parts.length);
    const partNames = parts.map(p => {
      return p.name
    })
    const indexRequestedPart = partNames.indexOf(requestedPart);

    if (!requestedPart || indexRequestedPart === -1) {
      console.log(requestedPart)
      console.log('Random section picked: ', `#${partIndex} - ${parts[partIndex].name}`);
    } else {
      partIndex = indexRequestedPart;
    }

    return Promise.all([
      parts[partIndex].name,
      api.fetchFromRosette(endpoint, parts[partIndex].unsafe),
    ]);
  })
  .then(results => {
    const name = results[0];
    const rosette = results[1];
    let formatted = {};

    console.log(`\n\n`);
    console.log(`${pageName} - ${name}`);
    console.log(`Endpoint: ${endpoint}`);

    if ( endpoint === 'topics' ) {
      console.log(`-> Key phrases`)
      if ( requestedOutput !== 'collectconcepts') {
        outputResult(requestedOutput, [pageName, name], rosette.keyphrases)
      }
      // console.table(rosette.keyphrases)
      console.log(`-> Concepts`)
      outputResult(requestedOutput, [pageName, name], rosette.concepts)
      // console.table(rosette.concepts)
    } else {
      formatted = rosette[endpoint]
        .map(ent => {
          delete ent.normalized
          delete ent.mentionOffsets
          return ent
        })
        .sort((a, b) => {
          // Sort by count, descending
          if (a.count < b.count) {
            return 1;
          } else if (a.count > b.count) {
            return -1;
          }
          return 0;
        })
        // console.table(formatted);
        outputResult(requestedOutput, [pageName, name], formatted)
      }
  })
  .catch(e => {
    console.log('Error:', (e.message || e))
    // console.log(e)
  })


function outputResult(outputType, pageNameParts, obj) {
  if ( outputType === 'csv' ) {
    const out = [];
    obj.forEach(res => {
      let row = Object.values(res);
      out.push(row.join(','))
    })
    console.log(out.join("\n"));
  } else if ( endpoint === 'topics' && outputType === 'collectconcepts' ) {
    // Collect terms per page-part
    // This is super specific for the experiment to collect concepts per pages
    // for the experiments folder. Ignore this for general usage ;)
    const out = {};
    const lim = Math.min(obj.length, 10);
    for ( let i = 0; i < lim; i++) {
      const res = obj[i]
      if (res.conceptId.indexOf('Q') === 0) {
        out[pageNameParts.join(' - ')] = out[pageNameParts.join(' - ')] || [];
        out[pageNameParts.join(' - ')].push({
          concept: `${res.phrase} (${res.conceptId})`,
          salience: res.salience
        });
      }
    }
    console.log(JSON.stringify(out, null, 2))
  } else {
    console.table(obj);
  }  
}