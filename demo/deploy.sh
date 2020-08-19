#!/usr/bin/env sh
# abort on errors
set -e
# build
echo "Building demo files..."
npm run build
# navigate into the build output directory
echo "Preparing to deploy..."
cd dist
# Start a repo
git init
git add -A
git commit -m 'deploy'
echo "Pushing branch..."
git push -f git@github.com:wikimedia/phoenix.git master:gh-pages
cd -
echo "Done."