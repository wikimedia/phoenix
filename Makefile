

build: 
	cd ./storage  && $(MAKE) build && cd ../
	cd ./common  && $(MAKE) build && cd ../
	cd ./lambdas/fetch-changed  && $(MAKE) build && cd ../../
	cd ./lambdas/fetch-schema.org  && $(MAKE) build && cd ../../
	cd ./lambdas/merge-schema.org  && $(MAKE) build && cd ../../
	cd ./lambdas/transform-parsoid  && $(MAKE) build && cd ../../

deploy: 
	cd ./lambdas/fetch-changed  && $(MAKE) deploy && cd ../../
	cd ./lambdas/fetch-schema.org  && $(MAKE) deploy && cd ../../
	cd ./lambdas/merge-schema.org  && $(MAKE) deploy && cd ../../
	cd ./lambdas/transform-parsoid  && $(MAKE) deploy && cd ../../	

test: 
	cd ./storage  && $(MAKE) test && cd ../
	cd ./common  && $(MAKE) test && cd ../
	cd ./lambdas/fetch-changed  && $(MAKE) test && cd ../../
	cd ./lambdas/fetch-schema.org  && $(MAKE) test && cd ../../
	cd ./lambdas/merge-schema.org  && $(MAKE) test && cd ../../
	cd ./lambdas/transform-parsoid  && $(MAKE) test && cd ../../
	