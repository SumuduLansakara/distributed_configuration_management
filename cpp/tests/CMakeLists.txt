find_package(Catch2 REQUIRED)

add_executable(tests main.cpp)
target_link_libraries(tests
        cfgstore
        Catch2::Catch2WithMain
        )

add_test(cctest tests)
