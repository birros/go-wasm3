build/wasm3:
	git clone --branch v0.5.0 --depth 1 https://github.com/wasm3/wasm3.git build/wasm3 && rm -rf build/wasm3/.git
	cd build/wasm3 && patch -p1 < ../../wasm3.patch

build/ios-cmake:
	git clone https://github.com/leetal/ios-cmake.git --branch 4.2.0 --depth 1 build/ios-cmake && rm -rf build/ios-cmake/.git

build/dockcross-windows-static-x64:
	mkdir -p build
	docker run --rm dockcross/windows-static-x64:20211018-b3b207e > build/dockcross-windows-static-x64
	chmod +x build/dockcross-windows-static-x64

lib/windows/amd64/libm3.a: build/wasm3 build/dockcross-windows-static-x64
	mkdir -p build/wasm3/build-windows-amd64
	cd build/wasm3 && ../dockcross-windows-static-x64 cmake -B build-windows-amd64
	cd build/wasm3 && ../dockcross-windows-static-x64 make -C build-windows-amd64
	mkdir -p lib/windows/amd64
	cp build/wasm3/build-windows-amd64/source/libm3.a lib/windows/amd64/libm3.a

lib/macos/amd64/libm3.a: build/wasm3
	mkdir -p build/wasm3/build-macos-amd64
	cd build/wasm3/build-macos-amd64 && cmake .. && make
	mkdir -p lib/macos/amd64
	cp build/wasm3/build-macos-amd64/source/libm3.a lib/macos/amd64/libm3.a

lib/linux/amd64/libm3.a: build/wasm3
	mkdir -p build/wasm3/build-linux-amd64
	cd build/wasm3/build-linux-amd64 && cmake .. && make
	mkdir -p lib/linux/amd64
	cp build/wasm3/build-linux-amd64/source/libm3.a lib/linux/amd64/libm3.a

lib/linux/arm64/libm3.a: build/wasm3
	mkdir -p build/wasm3/build-linux-arm64
	cd build/wasm3/build-linux-arm64 && cmake .. && make
	mkdir -p lib/linux/arm64
	cp build/wasm3/build-linux-arm64/source/libm3.a lib/linux/arm64/libm3.a

lib/ios/arm64/libm3.a: build/wasm3 build/ios-cmake
	mkdir -p build/wasm3/build-ios-arm64
	cd build/wasm3/build-ios-arm64 && cmake .. -DPLATFORM=OS64 -DCMAKE_TOOLCHAIN_FILE=../../ios-cmake/ios.toolchain.cmake && make
	mkdir -p lib/ios/arm64
	cp build/wasm3/build-ios-arm64/source/libm3.a lib/ios/arm64/libm3.a

lib/android/arm64/libm3.a: build/wasm3
	mkdir -p build/wasm3/build-android-arm64
	cd build/wasm3/build-android-arm64 && cmake .. -DANDROID_ABI=arm64-v8a -DANDROID_NATIVE_API_LEVEL=28 -DCMAKE_TOOLCHAIN_FILE=${ANDROID_NDK_HOME}/build/cmake/android.toolchain.cmake && make
	mkdir -p lib/android/arm64
	cp build/wasm3/build-android-arm64/source/libm3.a lib/android/arm64/libm3.a

.PHONY: clean
clean:
	git clean -xdf
