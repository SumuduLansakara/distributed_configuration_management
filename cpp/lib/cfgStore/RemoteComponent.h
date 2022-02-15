#pragma once

#include "Component.h"
#include <string>

namespace cfgStore {

    class RemoteComponent : public Component {
    public:
        RemoteComponent(std::string const& kind, std::string const& name, std::string const& instanceId)
            : Component(kind, name, instanceId) {
        }

        void setParam(std::string const& key, std::string const& val) override;
        std::string getParam(std::string const& key) override;
        void removeParam(std::string const& key) override;
    };

}// namespace cfgStore
