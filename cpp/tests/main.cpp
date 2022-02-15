#define CATCH_CONFIG_MAIN

#include <catch2/catch_all.hpp>
#include <etcd/Client.hpp>

TEST_CASE("Etcd") {
    etcd::Client etcd("http://127.0.0.1:2379");

    SECTION("Setting") {
        auto resTask = etcd.set("/test/key", "43");
        auto res = resTask.get();
        REQUIRE(res.is_ok());
    }

    SECTION("Getting") {
        auto resTask = etcd.get("/test/key");
        auto res = resTask.get();
        REQUIRE(res.is_ok());
    }

    SECTION("Removing") {
        auto resTask = etcd.rm("/test/key");
        auto res = resTask.get();
        REQUIRE(res.is_ok());
    }

    SECTION("Getting") {
        auto resTask = etcd.get("/test/key");
        auto res = resTask.get();
        REQUIRE(!res.is_ok());
    }
}
