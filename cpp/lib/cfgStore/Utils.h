#pragma once

#include <string>
#include <vector>

namespace cfgStore::utils {
    std::string generateUuidV4();

    std::vector<std::string> split(const std::string& s, char delim);
}// namespace cfgStore::utils
