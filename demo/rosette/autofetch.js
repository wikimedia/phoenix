import ApiHelper from './ApiHelper.js'
import { createRequire } from 'module';
const require = createRequire(import.meta.url);
const fs = require('fs')
const pages = {
  fruits: [
    'Apple',
    'Banana',
    'Pear',
    'Peach',
    'Apricot',
    'Cherry',
    'Pineapple',
    'Kiwifruit',
    'Mango',
    'Fig',
    'Watermelon',
    'Papaya',
    'Dragonfruit',
    'Orange (fruit)',
    'Plum',
    'Passionfruit',
    'Rambutan',
    'Lychee',
    'Lemon',
    'Blueberry',
    'Raspberry',
    'Grapefruit',
    'Lime',
    'Strawberry',
    'Mandarin orange'
  ],
  scientists: [
    'Albert Einstein',
    'Galileo Galilei',
    'Christiaan Huygens',
    'Isaac Newton',
    'Leonhard Euler',
    // 'Joseph Louis Lagrange',
    'Pierre-Simon Laplace',
    'Joseph Fourier',
    'James Clerk Maxwell',
    // 'Hendrik A. Lorentz',
    'Nikola Tesla',
    'Max Planck',
    'Emmy Noether',
    'Max Born',
    'Niels Bohr',
    'Erwin Schrödinger',
    'Louis de Broglie',
    'Satyendra Nath Bose',
    'Enrico Fermi',
    'Abdus Salam',
    'Charles Darwin',
    'Elizabeth Blackburn',
    'Lorna Casselton',
    'Rachel Carson',
    'Rosalind Franklin',
    'Jane Goodall',
    'Dorothy Hodgkin',
    'Shirley Ann Jackson',
    'Cynthia Kenyon',
    'Ada Lovelace'
  ]
}


const finalData = fetchForPages(
  pages.scientists,
  'scientists'
);

// FUNCTIONS
    
async function fetchForPages(inpPageList, filename = 'output') {
  const startTime = new Date();
  const resultData = {};
  const api = new ApiHelper();
  // Remove duplicates, just in case
  const pageList = Array.from(new Set(inpPageList));
  

  for (let n = 0; n < pageList.length; n++ ) {
    const pageName = pageList[n];
    try {
      log('GraphQL', `Fetching parts for "${pageName}"`)
      const parts = await api.fetchFromGraphQL(pageName);

      for (let p = 0; p < parts.length; p++ ) {
        const part = parts[p]
        const title = `${pageName} - ${part.name}`;
        const wait = await api.sleep(1000);
        
        log('Rosette', `Fetching rosette concepts for "${title}"`)
        try {
          const topics = await api.fetchFromRosette('topics', part.unsafe);
          const concepts = topics.concepts.sort((a, b) => {
            // Sort by count, descending
            if (a.salience < b.salience) {
              return 1;
            } else if (a.salience > b.salience) {
              return -1;
            }
            return 0;
          });
          // Limit the topics to some sensible number
          // In "real life" we should go more by salience
          const lim = Math.min(concepts.length, 20);


          for ( let i = 0; i < lim; i++) {
            const res = concepts[i]
            // Only include if the concept ID is wikidata
            if (res.conceptId.indexOf('Q') === 0) {
              // Create the group if not yet created
              resultData[title] = resultData[title] || [];
              // Output
              resultData[title].push({
                title: {
                  page: pageName,
                  part: part.name
                },
                concept: `${res.phrase} (${res.conceptId})`,
                salience: res.salience
              });
            }
          }
        
        } catch (e) {
          log(`Rosette API (${title})`, (e.message || e), 'error')
        }

      }
    } catch (e) {
      log('GraphQL Service', (e.message || e), 'error')
    }
  }

  log('Output', `Writing to file: data/${filename}-parts.json`)
  writeToFile(resultData, `data/${filename}-parts.json`)

  // Reorder by topics
  log('Output', 'Formatting by topic.')
  const resultsByTopics = {}
  Object.keys(resultData).forEach(partName => {
    resultData[partName].forEach(topicData => {
      // initialize as object if not yet initialized
      resultsByTopics[topicData.concept] = resultsByTopics[topicData.concept] || {};
      // Add page with salience result
      // Kinda betting here that there will never be two pages with the exact same salience
      // ... which is a fair assumption for an experiment
      resultsByTopics[topicData.concept][topicData.salience] = topicData.title;
    })
  })
  // // Sort each topic by salience
  // Object.keys(resultsByTopics).forEach(topic => {
  //   resultsByTopics[topic] = Object.entries(resultsByTopics[topic])
  //     // Sort by key, descending
  //     .sort(([key1,val1], [key2,val2]) => +key1 - +key2)
  //     // Transform back into an object
  //     .reduce((r, [key, val]) => ({ ...r, [key]: val }), {});
  // })

  log('Output', `Writing to file: data/${filename}-topics.json`)
  writeToFile(resultsByTopics, `data/${filename}-topics.json`)

  console.log(`\n\n`)
  const endTime = new Date();
  log(`Success!`, `Data available at data/${filename}-parts.json and data/${filename}-topics.json`)
  console.log(`\n`)
  log(`Number of pages`, pageList.length)
  log(`Started at`, startTime.toTimeString())
  log(`Ended at`, endTime.toTimeString())
  console.log(`\n\n`)
}

function writeToFile(content, path) {
  try {
    fs.writeFileSync(path, JSON.stringify(content, null, 2))
  } catch (e) {
    log(`Error writing to file at ${path}`, e)
  }
}
/**
 * Output to the terminal
 *
 * @param {string} topic 
 * @param {string} str 
 */
function log(topic, str, type = 'log') {
  console.log(`[${type.toUpperCase()}] ${topic}: ${str}`)
}
