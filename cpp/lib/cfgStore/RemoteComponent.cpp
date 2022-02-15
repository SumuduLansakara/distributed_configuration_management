#include "RemoteComponent.h"
#include "Communicator.h"

namespace cfgStore {

    void RemoteComponent::setParam(std::string const& key, std::string const& val) {
        Communicator::set(paramPath(_instanceId, key), val);
    }

    std::string RemoteComponent::getParam(std::string const& key) {
        return Communicator::get(paramPath(_instanceId, key));
    }

    void RemoteComponent::removeParam(std::string const& key) {
        Communicator::rm(paramPath(_instanceId, key));
    }

}// namespace cfgStore
