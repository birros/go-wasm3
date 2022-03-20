WASM3_COMMIT=6b8bcb1e07bf26ebef09a7211b0a37a446eafd52
DOCKCROSS_COMMIT=9764a2ece1ec7913d62da981451f04acd3cc61cc
DOCKCROSS_TAG=20220308-9764a2e
DOCKCROSS_ORG=local-dockcross
ANDROID_NDK_REVISION=24
ANDROID_NDK_API=28

MAKEFILE_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build/wasm3:
	git clone https://github.com/wasm3/wasm3.git build/wasm3
	cd build/wasm3 && git checkout $(WASM3_COMMIT)
	cd build/wasm3 && patch -p1 < ../../patches/wasm3.patch

build/dockcross:
	git clone https://github.com/dockcross/dockcross.git build/dockcross
	cd build/dockcross && git checkout $(DOCKCROSS_COMMIT)
	cd build/dockcross && patch -p1 < ../../patches/dockcross.patch

build/ios-cmake:
	git clone https://github.com/leetal/ios-cmake.git --branch 4.2.0 --depth 1 build/ios-cmake && rm -rf build/ios-cmake/.git

build/dockcross-windows-static-x64:
	mkdir -p build
	docker run --rm dockcross/windows-static-x64:$(DOCKCROSS_TAG) > build/dockcross-windows-static-x64
	chmod +x build/dockcross-windows-static-x64

build/dockcross-linux-arm64:
	mkdir -p build
	docker run --rm dockcross/linux-arm64:$(DOCKCROSS_TAG) > build/dockcross-linux-arm64
	chmod +x build/dockcross-linux-arm64

build/dockcross-linux-x64:
	mkdir -p build
	docker run --rm dockcross/linux-x64:$(DOCKCROSS_TAG) > build/dockcross-linux-x64
	chmod +x build/dockcross-linux-x64

build/dockcross-android-arm64: build/dockcross
	mkdir -p build
    # patch dockcross/android-arm64 to change ANDROID_NDK_REVISION and ANDROID_NDK_API
	cd build/dockcross/android-arm64 && sed -i.bak 's/base:.*/base:$(DOCKCROSS_TAG)/g'                                        Dockerfile.in && rm Dockerfile.in.bak
	cd build/dockcross/android-arm64 && sed -i.bak 's/ANDROID_NDK_REVISION .*/ANDROID_NDK_REVISION $(ANDROID_NDK_REVISION)/g' Dockerfile.in && rm Dockerfile.in.bak
	cd build/dockcross/android-arm64 && sed -i.bak 's/ANDROID_NDK_API .*/ANDROID_NDK_API $(ANDROID_NDK_API)/g'                Dockerfile.in && rm Dockerfile.in.bak
	cd build/dockcross && make ORG=$(DOCKCROSS_ORG) android-arm64
    # end patch
	docker run --rm $(DOCKCROSS_ORG)/android-arm64 > build/dockcross-android-arm64
	chmod +x build/dockcross-android-arm64

lib/windows/amd64/libm3.a: build/wasm3 build/dockcross-windows-static-x64
	mkdir -p build/wasm3/build-windows-amd64
	cd build/wasm3 && ../dockcross-windows-static-x64 cmake -B build-windows-amd64
	cd build/wasm3 && ../dockcross-windows-static-x64 make -C build-windows-amd64
	mkdir -p lib/windows/amd64
	cp build/wasm3/build-windows-amd64/source/libm3.a lib/windows/amd64/libm3.a

lib/macos/amd64/libm3.a: build/wasm3
	mkdir -p build/wasm3/build-macos-amd64
	cd build/wasm3 && cmake -B build-macos-amd64
	cd build/wasm3 && make -C build-macos-amd64
	mkdir -p lib/macos/amd64
	cp build/wasm3/build-macos-amd64/source/libm3.a lib/macos/amd64/libm3.a

lib/linux/amd64/libm3.a: build/wasm3 build/dockcross-linux-x64
	mkdir -p build/wasm3/build-linux-amd64
	cd build/wasm3 && ../dockcross-linux-x64 cmake -B build-linux-amd64
	cd build/wasm3 && ../dockcross-linux-x64 make -C build-linux-amd64
	mkdir -p lib/linux/amd64
	cp build/wasm3/build-linux-amd64/source/libm3.a lib/linux/amd64/libm3.a

lib/linux/arm64/libm3.a: build/wasm3 build/dockcross-linux-arm64
	mkdir -p build/wasm3/build-linux-arm64
	cd build/wasm3 && ../dockcross-linux-arm64 cmake -B build-linux-arm64
	cd build/wasm3 && ../dockcross-linux-arm64 make -C build-linux-arm64
	mkdir -p lib/linux/arm64
	cp build/wasm3/build-linux-arm64/source/libm3.a lib/linux/arm64/libm3.a

lib/ios/arm64/libm3.a: build/wasm3 build/ios-cmake
	mkdir -p build/wasm3/build-ios-arm64
	cd build/wasm3 && cmake -B build-ios-arm64 -DPLATFORM=OS64 -DCMAKE_TOOLCHAIN_FILE=$(MAKEFILE_DIR)build/ios-cmake/ios.toolchain.cmake
	cd build/wasm3 && make -C build-ios-arm64
	mkdir -p lib/ios/arm64
	cp build/wasm3/build-ios-arm64/source/libm3.a lib/ios/arm64/libm3.a

lib/android/arm64/libm3.a: build/wasm3 build/dockcross-android-arm64
	mkdir -p build/wasm3/build-android-arm64
	cd build/wasm3 && ../dockcross-android-arm64 cmake -B build-android-arm64
	cd build/wasm3 && ../dockcross-android-arm64 make -C build-android-arm64
	mkdir -p lib/android/arm64
	cp build/wasm3/build-android-arm64/source/libm3.a lib/android/arm64/libm3.a

.PHONY: clean
clean:
	rm -rf build
	make -C examples/sum clean
