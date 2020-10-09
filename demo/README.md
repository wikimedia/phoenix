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
4. Run `npm run serve` 
5. In another terminal window, run `npm run phoenix` to activate the GraphQL service

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

## Deploy to Github pages
To deploy the demo to github pages:

1. Clone the Phoenix repo
2. Go into the `phoenix/demo` folder
3. Run `./deploy.sh`

The demo will be avilable at the [Wikimedia Github Pages for the phoenix repo](https://wikimedia.github.io/phoenix).

NOTE: Running that command will build the repository and then immediately push this into github. Make sure your working directory is clean.

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
