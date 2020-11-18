build:
	cd lib && gcc -fPIC -c greengrasssdk.c && \
	gcc -shared -o ../libaws-greengrass-core-sdk-c.so greengrasssdk.o && \
	rm -f greengrasssdk.o && cd ..

clean:
	@rm -rf lib/*.o *.so