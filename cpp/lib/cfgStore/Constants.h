#pragma once

namespace cfgStore {
    constexpr auto ETCD_CONN_STR = "http://127.0.0.1:2379";

    constexpr auto PREFIX_COMPONENTS = "/components";
    constexpr auto PREFIX_PARAMETERS = "/parameters";
    constexpr auto PREFIX_METADATA = "/meta";
}// namespace cfgStore
