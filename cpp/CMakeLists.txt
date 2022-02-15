project(app)
cmake_minimum_required(VERSION 3.15)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
set(CMAKE_CXX_STANDARD 20)

add_subdirectory(lib)
add_subdirectory(demo)

enable_testing()
add_subdirectory(tests)
