build:
	@mkdir -p /tmp/gogreengrass
	@cd lib && gcc -fPIC -c greengrasssdk.c && \
	gcc -shared -o /tmp/gogreengrass/libaws-greengrass-core-sdk-c.so greengrasssdk.o && \
	rm -f greengrasssdk.o && cd ..

clean:
	@rm -rf lib/*.o /tmp/gogreengrass/libaws-greengrass-core-sdk-c.so