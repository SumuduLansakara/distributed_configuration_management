#pragma once

#include "RemoteComponent.h"
#include <string>
#include <utility>

namespace cfgStore {
    class RemoteComponentStub {
    public:
        RemoteComponentStub(std::string kind, std::string instanceId)
            : _kind{std::move(kind)}
            , _instanceId{std::move(instanceId)} {
        }

        std::string const& getKind() { return _kind; }
        std::string const& getInstanceId() { return _instanceId; }

        RemoteComponent spawn();

    private:
        std::string const _kind;
        std::string const _instanceId;
    };
}// namespace cfgStore
