
test: 
	cd ./storage  && $(MAKE) test && cd ../
	cd ./common  && $(MAKE) test && cd ../
	cd ./lambdas/fetch-changed  && $(MAKE) test && cd ../../
	cd ./lambdas/fetch-schema.org  && $(MAKE) test && cd ../../
	cd ./lambdas/merge-schema.org  && $(MAKE) test && cd ../../
	cd ./lambdas/transform-parsoid  && $(MAKE) test && cd ../../
	