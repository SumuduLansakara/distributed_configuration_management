find_package(OpenSSL REQUIRED)
find_package(Protobuf REQUIRED)
find_library(CPPREST_LIB cpprest)
find_package(fmt REQUIRED)

include(${PROJECT_SOURCE_DIR}/cmake/FindGRPC.cmake)
set(GRPC_LIBRARIES ${GPR_LIBRARY} ${GRPC_LIBRARY} ${GRPC_GRPC++_LIBRARY})

add_library(cfgstore
        cfgStore/LocalComponent.cpp
        cfgStore/RemoteComponentStub.cpp
        cfgStore/RemoteComponent.cpp
        cfgStore/Utils.cpp
        )
target_link_libraries(cfgstore PRIVATE
        ${CPPREST_LIB}
        etcd-cpp-api
        OpenSSL::SSL
        fmt::fmt
        ${Protobuf_LIBRARIES}
        ${GRPC_LIBRARIES}
        )
target_include_directories(cfgstore PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR})
