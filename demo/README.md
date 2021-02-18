# Phoenix demo

This sets up an interactive website that demos the structured content store capabilities with GraphQL.

## Setup

To run this locally:

1. Clone the Phoenix repo
2. Go into the `phoenix/demo` folder
3. Run `npm install`
4. Run `npm run build`
5. The deployable site will be available in the `/dist` folder

NOTE: If running this locally, please remember to activate the `phoenix/service` service for the interface to have the GraphQL API access.

## Development

To develop and work on this code locally:

1. Clone the Phoenix repo
2. Go into the `phoenix/demo` folder
3. Run `npm install`
4. Run `npm run serve` to run against the published GraphQL AWS endpoint

## Development against local GraphQL service

If you intend to test locally against a local instance of the GraphQL service, follow steps 1-3 above, and then:

4. Run `npm run local` to run against the local GraphQL endpoint
5. In another terminal window, run `npm run phoenix` to activate the GraphQL service locally

### Development tools

Lints and fixes files
```
npm run lint
```

Compiles and hot-reloads for development
```
npm run serve
```

Compiles and minifies for production
```
npm run build
```

Production-ready files will be available in `/dist`

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
