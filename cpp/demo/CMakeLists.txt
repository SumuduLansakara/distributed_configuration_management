add_executable(demo
        Parent.cpp
        Child.cpp
        main.cpp
        )

target_link_libraries(demo PRIVATE
        cfgstore
        )

target_include_directories(demo PRIVATE
        /usr/local/include/etcd/proto
        )
