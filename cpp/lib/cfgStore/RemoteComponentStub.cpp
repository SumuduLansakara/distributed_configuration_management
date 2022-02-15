#include "RemoteComponentStub.h"
#include "Communicator.h"

namespace cfgStore {
    RemoteComponent RemoteComponentStub::spawn() {
        auto const metaPrefix = metadataPath(_instanceId);
        auto const metaRes = Communicator::ls(metaPrefix);
        std::map<std::string, std::string> metaData;
        for (int i = 0; i < metaRes.keys().size(); ++i) {
            auto const& key = metaRes.key(i).substr(metaPrefix.size());
            auto const& val = metaRes.value(i).as_string();
            metaData.emplace(key, val);
        }
        return {metaData["kind"], metaData["name"], _instanceId};
    }
}// namespace cfgStore
