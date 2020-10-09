import axios from 'axios'
import { createRequire } from 'module';
const require = createRequire(import.meta.url);
const { ArgumentParser } = require('argparse');

import Api from 'rosette-api';
import config from './config.js'

var parser = new ArgumentParser({
  add_help: true,
  description: "Fetch data for a section part from rosette"
});

parser.add_argument("--page", {help: "Page name from simple.wikipedia", required: false});
parser.add_argument("--part", {help: "A section name to fetch", required: false});
parser.add_argument("--endpoint", {help: "Endpoint for Rosette API", required: false});
const args = parser.parse_args();
const rosetteApi = new Api(config.rosette_key, args.url);

// Available endpoints: 'topics', 'entities', 'categories'
const endpoint = args.endpoint || "topics";
const pageName = args.page || 'Philadelphia'
const requestedPart = args.part || ''
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
axios
  .post(
    'http://localhost:8080/query',
    { query: queries.parts },
    { headers: { 'Content-Type': 'application/json' } }
  )
  .then(res => {
    if ( res && res.data && res.data.data && ( res.data.data.page || res.data.data.node ) ) {
      return res.data.data.page.hasPart
    }
    return Promise.reject('Requested node or page not found.');
  })
  .then(parts => {
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
      fetchFromRosette(endpoint, parts[partIndex].unsafe),
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
      console.table(rosette.keyphrases)
      console.log(`-> Concepts`)
      console.table(rosette.concepts)
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
        console.table(formatted);
    }
  })
  .catch(e => {
    console.log('Error:', (e.message || e))
    // console.log(e)
  })


function fetchFromRosette(api_endpoint, content) {
  const deferred = new Deferred();
  rosetteApi.parameters.content = content;
  rosetteApi.rosette(api_endpoint, function(err, res){
    if(err){
      deferred.reject(err);
    } else {
      deferred.resolve(res);
    }
  });
  return deferred;
}

function Deferred () {
  var res = null,
    rej = null,
    p = new Promise((a,b)=>(res = a, rej = b));
  p.resolve = res;
  p.reject = rej;
  return p;
}
