# Wikipedia Phoenix - Rosette experiment

This is a script to enable an experimentation of the Rosette API drawing from the Phoenix storage content.

## Usage

To run the script:

1. Clone the repo
2. Go to `demo/rosette`
3. Run `npm install`
4. Sign up for Rosette API key.
5. Rename `config.sample.js` to `config.js` and fill in our Rosette API key.

### Using the script

In order to use this script, our GraphQL service needs to run locally:
1. Go to `service/`
2. Follow the README instructions to run the GraphQL service.
3. Run the script in `demo/rosette/index.js` using the parameters below.

#### Available parameters

* `--page "[name]"` Name of the requested page from Phoenix storage. Defaults to "Philadelphia"
* `--part "[name]"` Name of the section from the requested page. If not given, a random section from the given page is used.
* `--endpoint "[name]"` Allows the change the endpoint in Rosette. Available endpoints: `topics`, `entities`, `categories`, `relationships`. If none given, the default is `topics` 

## Some examples

### Categories:
Potentially too generic for our use? 

* `node index.js --endpoint "categories" --page "Yellow (song)" --part "Composition"`
* Empty - `node index.js --endpoint "categories" --page "Philadelphia" --part "History"`
* Generic - `node index.js --endpoint "categories" --page "Philadelphia" --part "Images"`

### Entities
Not sure where to use this unless this is used for indexing/search.

* `node index.js --endpoint "entities" --page "Yellow (song)" --part "__intro__"`
* Generic - `node index.js --endpoint "entities" --page "Philadelphia" --part "History"`

### Topics:
Probably most useful for us.

* `node index.js --endpoint "topics" --page "Yellow (song)" --part "Release and reception"`

### Relationships:
Potentially useful for us.

See [Relationship types](https://developer.rosette.com/features-and-functions#relationship-extraction-relationship-types)

* `node index.js --endpoint "relationships" --page "Philadelphia" --part "Culture"`
* Empty - `node index.js --endpoint "relationships" --page "Banana" --part "__intro__"`

# Automated script

The automated script, `autofetch.js` is meant to pick a series of pages and pass their individual parts through Rosette's `topics` endpoint, gathering and mapping the sections per topic.

The list of pages are defined inside the `autofetch.js` file.

The results are created inside the `data/` folder.