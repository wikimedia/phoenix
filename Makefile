


clean:
	$(MAKE) -c storage clean
	$(MAKE) -c common clean
	$(MAKE) -c fetch-changed clean
	$(MAKE) -c fetch-schema.org clean
	$(MAKE) -c merge-schema.org clean
	$(MAKE) -c transform-parsoid clean

build: 
	$(MAKE) -C storage build
	$(MAKE) -C common build
	$(MAKE) -C lambdas/fetch-changed build
	$(MAKE) -C lambdas/fetch-schema.org  build
	$(MAKE) -C lambdas/merge-schema.org  build
	$(MAKE) -C lambdas/transform-parsoid build

deploy: 
	$(MAKE) -C lambdas/fetch-changed deploy
	$(MAKE) -C lambdas/fetch-schema.org  deploy
	$(MAKE) -C lambdas/merge-schema.org  deploy
	$(MAKE) -C lambdas/transform-parsoid deploy

test: 
	$(MAKE) -C storage test
	$(MAKE) -C common test
	$(MAKE) -C lambdas/fetch-changed test
	$(MAKE) -C lambdas/fetch-schema.org  test
	$(MAKE) -C lambdas/merge-schema.org  test
	$(MAKE) -C lambdas/transform-parsoid test

.PHONY: build clean deploy test
